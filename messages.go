package dhcp

import (
	"fmt"
	"net"
)

// Message represents a DHCP
// message
type Message struct {
	Type          byte
	TransactionID uint32
	ServerIP      net.IP
	ClientIP      net.IP
	ClientMAC     net.HardwareAddr

	Options Options
}

// NewMessage creates a new DHCP message with
// the specified parameters
func NewMessage(t byte, transactionId uint32, serverIp, clientIp net.IP, clientMac net.HardwareAddr) Message {
	var m Message
	m.Type = t
	m.TransactionID = transactionId
	m.ServerIP = serverIp
	m.ClientIP = clientIp
	m.ClientMAC = clientMac
	m.Options = make(map[byte][]byte)

	return m
}

// MessageFromFrame retreives the values from a DHCP frame
// and constructs a Message from it
func MessageFromFrame(f frame) (Message, error) {
	var m Message

	t, ok := f.opts[OptionDHCPMessageType]
	if !ok {
		return m, fmt.Errorf("no dhcp message type in frame")
	}

	m.Type = t[0]
	m.TransactionID = f.xid
	m.ServerIP = net.IP(f.siaddr).To4()
	m.ClientMAC = net.HardwareAddr(f.chaddr[:f.hlen])

	if m.Type == DHCPTypeRequest {
		ip, ok := f.opts[OptionRequestedIPAddress]
		if !ok {
			return m, fmt.Errorf("no requested ip address in dhcp request")
		}

		m.ClientIP = net.IP(ip).To4()
	}

	return m, nil
}

func (m *Message) SetOption(option byte, value []byte) {
	m.Options[option] = value
}

func (m Message) GetFrame() []byte {
	var f frame

	if m.Type == DHCPTypeDiscover || m.Type == DHCPTypeRequest || m.Type == DHCPTypeRelease {
		f.op = 0x01
	}
	if m.Type == DHCPTypeOffer || m.Type == DHCPTypeACK || m.Type == DHCPTypeNACK || m.Type == DHCPTypeDecline {
		f.op = 0x02
	}

	f.htype = 0x01
	f.hlen = 0x06
	f.hops = 0x00
	f.xid = m.TransactionID
	f.secs = 0x0000
	f.flags = 0x0000
	f.ciaddr = unpack(4, uint64(0x00000000))
	f.yiaddr = m.ClientIP
	f.siaddr = m.ServerIP
	f.giaddr = unpack(4, uint64(0x00000000))
	f.chaddr = m.ClientMAC

	f.opts = make(map[byte][]byte)

	for opt, val := range m.Options {
		f.opts[opt] = val
	}

	f.opts[OptionDHCPMessageType] = []byte{m.Type}

	return f.toBytes()
}
