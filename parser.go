package dhcp

import (
	"fmt"
	"net"
)

type frame struct {
	src  net.Addr
	size int

	hlen   uint8
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
	cookie []byte

	opts map[byte][]byte
}

func parse(socket net.PacketConn, buf []byte) (frame, error) {
	var f frame
	f.opts = make(map[byte][]byte)

	n, addr, err := socket.ReadFrom(buf)
	if err != nil {
		return f, err
	}
	if n <= 236 {
		return f, fmt.Errorf("invalid dhcp packet")
	}

	f.src = addr
	f.size = n

	f.hlen = uint8(buf[2])
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
	f.cookie = buf[236:240]

	cursor := 0
	optbuf := buf[240:]

	for cursor < len(optbuf) {
		code := optbuf[cursor]
		if code == OptionPad {
			continue
		}
		if code == OptionEnd {
			break
		}

		len := optbuf[cursor+1]
		cursor += 2

		f.opts[code] = optbuf[cursor : cursor+int(len)]
		cursor += int(len)
	}

	return f, nil
}
