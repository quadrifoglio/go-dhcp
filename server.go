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
	socket *net.UDPConn

	discoverCb DiscoverCallback
	requestCb  RequestCallback
	releaseCb  ReleaseCallback
}

func NewServer() (Server, error) {
	var serv Server
	serv.discoverCb = nil
	serv.requestCb = nil
	serv.releaseCb = nil

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
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 67})
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
			if s.requestCb != nil {
				var reqIp net.IP

				if ip, ok := frame.opts[OptionRequestedIPAddress]; ok {
					reqIp = ip
				} else {
					return fmt.Errorf("no requested ip address in dhcp request")
				}

				s.requestCb(s, frame.xid, frame.chaddr[:frame.hlen], reqIp)
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
