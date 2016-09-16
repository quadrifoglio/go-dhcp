package dhcp

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	discoverCb func(net.PacketConn, net.HardwareAddr)
	requestCb  func(net.PacketConn, net.HardwareAddr, net.IP)
	releaseCb  func(net.PacketConn, net.HardwareAddr)
}

func NewServer() (Server, error) {
	var serv Server
	return serv, nil
}

func (s *Server) HandleDiscover(callback func(net.PacketConn, net.HardwareAddr)) {
	s.discoverCb = callback
}

func (s *Server) HandleRequest(callback func(net.PacketConn, net.HardwareAddr, net.IP)) {
	s.requestCb = callback
}

func (s *Server) HandleRelease(callback func(net.PacketConn, net.HardwareAddr)) {
	s.releaseCb = callback
}

func (s *Server) Start() error {
	socket, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		return err
	}

	defer socket.Close()

	return s.run(socket)
}

func (s *Server) run(socket net.PacketConn) error {
	buf := make([]byte, 1500)

	for {
		n, _, err := socket.ReadFrom(buf)
		if err != nil {
			return err
		}

		frame, err := parse(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}

		var msgType byte

		if b, ok := frame.opts[OptionDHCPMessageType]; ok {
			msgType = b[0]
		} else {
			return fmt.Errorf("no valid dhcp message type field")
		}

		if msgType == DHCPTypeDiscover {
			if s.discoverCb != nil {
				s.discoverCb(socket, frame.chaddr[:frame.hlen])
			}
		}
		if msgType == DHCPTypeRequest {
			var requestIP net.IP

			if ipBytes, ok := frame.opts[OptionRequestedIPAddress]; ok {
				requestIP = net.IP(ipBytes)
			} else {
				return fmt.Errorf("no request ip address in dhcp request")
			}

			/*if srvBytes, ok := frame.opts[OptionServerIdentifier]; ok {
				info.RequestDHCPServer = net.IP(srvBytes)
			} else {
				continue
			}*/

			if s.requestCb != nil {
				s.requestCb(socket, frame.chaddr[:frame.hlen], requestIP)
			}
		}
		if msgType == DHCPTypeRelease {
			if s.releaseCb != nil {
				s.releaseCb(socket, frame.chaddr[:frame.hlen])
			}
		}
	}

	return nil
}
