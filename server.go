package dhcp

import (
	"log"
	"net"
)

type Server struct {
}

func NewServer(startAddr string, numLease int, opts []Option) (Server, error) {
	var serv Server

	return serv, nil
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
		_, err := parse(socket, buf)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}
