package dhcp

import (
	"fmt"
	"log"
	"net"
)

type DiscoverCallback func(*Server, uint32, net.HardwareAddr)
type RequestCallback func(*Server, uint32, net.HardwareAddr, net.IP)
type ReleaseCallback func(*Server, net.HardwareAddr)

type Server struct {
	socket net.PacketConn

	discoverCb DiscoverCallback
	requestCb  RequestCallback
	releaseCb  ReleaseCallback
}

func NewServer() (Server, error) {
	var serv Server
	return serv, nil
}

func (s *Server) BroadcastPacket(packet []byte) error {
	addr := &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 68}

	_, err := s.socket.WriteTo(packet, addr)
	if err != nil {
		return err
	}

	return nil
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

func (s *Server) ListenAndServe() error {
	socket, err := net.ListenPacket("udp4", ":67")
	if err != nil {
		return err
	}

	s.socket = socket
	defer s.socket.Close()

	buf := make([]byte, 1500)

	for {
		n, _, err := s.socket.ReadFrom(buf)
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
				s.discoverCb(s, frame.xid, frame.chaddr[:frame.hlen])
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
				s.requestCb(s, frame.xid, frame.chaddr[:frame.hlen], requestIP)
			}
		}
		if msgType == DHCPTypeRelease {
			if s.releaseCb != nil {
				s.releaseCb(s, frame.chaddr[:frame.hlen])
			}
		}
	}

	return nil
}
