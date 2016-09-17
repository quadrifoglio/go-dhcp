package dhcp

import (
	"fmt"
	"log"
	"net"
)

type DiscoverCallback func(net.PacketConn, uint32, net.HardwareAddr)
type RequestCallback func(net.PacketConn, uint32, net.HardwareAddr, net.IP)
type ReleaseCallback func(net.PacketConn, net.HardwareAddr)

type Server struct {
	discoverCb DiscoverCallback
	requestCb  RequestCallback
	releaseCb  ReleaseCallback
}

func NewServer() (Server, error) {
	var serv Server
	return serv, nil
}

func (s *Server) HandleDiscover(callback DiscoverCallback) {
	s.discoverCb = callback
}

func (s *Server) HandleRequest(callback RequestCallback) {
	s.requestCb = callback
}

func (s *Server) HandleRelease(callback ReleaseCallback) {
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
				s.discoverCb(socket, frame.xid, frame.chaddr[:frame.hlen])
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
				s.requestCb(socket, frame.xid, frame.chaddr[:frame.hlen], requestIP)
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
