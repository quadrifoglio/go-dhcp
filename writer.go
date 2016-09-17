package dhcp

import (
	"net"
)

type Offer struct {
	serverIP    net.IP
	clientMAC   []byte
	transaction uint32

	LeaseTime int
	IP        net.IP
	Mask      net.IP
	Router    net.IP
	DNS       []net.IP
}

func NewOffer(serverIP net.IP, clientMAC []byte, transaction uint32) Offer {
	var o Offer
	o.serverIP = serverIP
	o.clientMAC = clientMAC
	o.transaction = transaction

	return o
}

func (o Offer) GetBytes() []byte {
	var f frame
	f.op = 0x02
	f.htype = 0x01
	f.hlen = 0x06
	f.hops = 0x00
	f.xid = o.transaction
	f.secs = 0x0000
	f.flags = 0x0000
	f.ciaddr = unpack(4, uint64(0x00000000))
	f.yiaddr = o.IP
	f.siaddr = o.serverIP
	f.giaddr = unpack(4, uint64(0x00000000))
	f.chaddr = o.clientMAC

	f.opts = make(map[byte][]byte)
	f.opts[OptionServerIdentifier] = o.serverIP.To4()
	f.opts[OptionDHCPMessageType] = []byte{DHCPTypeOffer}
	f.opts[OptionSubnetMask] = o.Mask.To4()
	f.opts[OptionIPAddressLeaseTime] = []byte{0x00, 0x00, 0xff, 0xff}

	if len(o.Router) > 0 {
		f.opts[OptionRouter] = o.Router
	}

	// TODO: Other options

	return f.toBytes()
}
