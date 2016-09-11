package dhcp

import (
	"fmt"
	"log"
	"net"
)

type Info struct {
	MAC []byte

	// DHCPRequest
	RequestIP         net.IP
	RequestDHCPServer net.IP
}

type Callback func(Info)

type Server struct {
	discoverCb Callback
	requestCb  Callback
	releaseCb  Callback
}

func NewServer(startAddr string, numLease int) (Server, error) {
	var serv Server

	return serv, nil
}

func (s *Server) HandleFunc(event string, callback Callback) {
	if event == "discover" {
		s.discoverCb = callback
	}
	if event == "request" {
		s.requestCb = callback
	}
	if event == "release" {
		s.releaseCb = callback
	}
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
		frame, err := parse(socket, buf)
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
			var info Info
			info.MAC = frame.chaddr[:frame.hlen]

			if s.discoverCb != nil {
				s.discoverCb(info)
			}
		}
		if msgType == DHCPTypeRequest {
			var info Info
			info.MAC = frame.chaddr[:frame.hlen]

			if ipBytes, ok := frame.opts[OptionRequestedIPAddress]; ok {
				info.RequestIP = net.IP(ipBytes)
			} else {
				continue
			}

			if srvBytes, ok := frame.opts[OptionServerIdentifier]; ok {
				info.RequestDHCPServer = net.IP(srvBytes)
			} else {
				continue
			}

			if s.requestCb != nil {
				s.requestCb(info)
			}
		}
		if msgType == DHCPTypeRelease {
			var info Info
			info.MAC = frame.chaddr[:frame.hlen]

			if s.releaseCb != nil {
				s.releaseCb(info)
			}
		}
	}

	return nil
}
