package dhcp

import (
	"fmt"
)

type frame struct {
	op     byte
	htype  byte
	hlen   byte
	hops   byte
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

func parse(buf []byte) (frame, error) {
	var f frame
	f.opts = make(map[byte][]byte)

	if len(buf) <= 236 {
		return f, fmt.Errorf("invalid dhcp packet")
	}

	f.op = buf[0]
	f.htype = buf[1]
	f.hlen = buf[2]
	f.hops = buf[3]
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

func (f frame) toBytes() []byte {
	buf := make([]byte, 240)

	buf[0] = f.op
	buf[1] = f.htype
	buf[2] = f.hlen
	buf[3] = f.hops

	copy(buf[4:], unpack(4, uint64(f.xid)))
	copy(buf[8:], unpack(2, uint64(f.secs)))
	copy(buf[10:], unpack(2, uint64(f.flags)))
	copy(buf[12:], f.ciaddr[:4])
	copy(buf[16:], f.yiaddr[:4])
	copy(buf[20:], f.siaddr[:4])
	copy(buf[24:], f.giaddr[:4])
	copy(buf[28:], f.chaddr[:f.hlen])

	// BOOTP legacy...
	var i int
	for i = 44; i < 44+192; i++ {
		buf[i] = 0
	}

	copy(buf[i:], []byte{0x63, 0x82, 0x53, 0x63})
	i += 4

	for code, value := range f.opts {
		buf = append(buf, code)
		buf = append(buf, byte(len(value)))
		buf = append(buf, value...)
	}

	buf = append(buf, OptionEnd)
	return buf
}
