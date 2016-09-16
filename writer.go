package dhcp

import (
	"net"
)

type Offer struct {
	serverIP  net.IP
	clientMAC []byte

	LeaseTime int
	IP        net.IP
	Mask      net.IP
	Router    net.IP
	DNS       []net.IP
}

func NewOffer(serverIP net.IP, clientMAC []byte) Offer {
	var o Offer
	o.serverIP = serverIP
	o.clientMAC = clientMAC

	return o
}

func (o Offer) GetBytes() []byte {
	var f frame
	f.op = 0x02
	f.htype = 0x01
	f.hlen = 0x06
	f.hops = 0x00
	f.xid = 0xabcdef01 // TODO: Random number
	f.secs = 0x0000
	f.flags = 0x0000
	f.ciaddr = unpack(4, uint64(0x00000000))
	f.yiaddr = o.IP
	f.siaddr = o.serverIP
	f.giaddr = unpack(4, uint64(0x00000000))
	f.chaddr = o.clientMAC

	f.opts = make(map[byte][]byte)
	f.opts[OptionDHCPMessageType] = []byte{DHCPTypeOffer}
	f.opts[OptionSubnetMask] = o.Mask

	// TODO: Other options

	return f.toBytes()
}
