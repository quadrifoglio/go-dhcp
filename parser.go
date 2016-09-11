package dhcp

import (
	"fmt"
	"net"
)

type frame struct {
	src  net.Addr
	size int

	xid    uint32
	secs   uint16
	flags  uint16
	ciaddr []byte
	yiaddr []byte
	siaddr []byte
	giaddr []byte
	chaddr []byte
	sname  []byte
	file   []byte
}

func parse(socket net.PacketConn, buf []byte) (frame, error) {
	var f frame

	n, addr, err := socket.ReadFrom(buf)
	if err != nil {
		return f, err
	}
	if n <= 236 {
		return f, fmt.Errorf("invalid dhcp packet")
	}

	f.src = addr
	f.size = n

	f.xid = uint32(pack(4, buf[4:8]))
	f.secs = uint16(pack(2, buf[8:10]))
	f.flags = uint16(pack(2, buf[10:12]))
	f.ciaddr = buf[12:16]
	f.yiaddr = buf[16:20]
	f.siaddr = buf[20:24]
	f.giaddr = buf[24:28]
	f.chaddr = buf[28:44]
	f.sname = buf[44:108]
	f.file = buf[108:236]

	return f, nil
}
