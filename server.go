package dhcp

import (
	"log"
	"net"
)

type Callback func(*Server, Message)

type Server struct {
	socket *net.UDPConn
	cb     Callback
}

func NewServer() (Server, error) {
	var serv Server
	serv.cb = nil

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

func (s *Server) HandleFunc(callback Callback) {
	s.cb = callback
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

		msg, err := MessageFromFrame(frame)
		if err != nil {
			log.Println(err)
			continue
		}

		s.cb(s, msg)
	}

	return nil
}
