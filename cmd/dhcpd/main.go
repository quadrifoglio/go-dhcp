package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/quadrifoglio/go-dhcp"
)

var (
	ServerIP = []byte{10, 0, 0, 1}

	LeaseIP     = []byte{10, 0, 0, 100}
	LeaseMask   = []byte{255, 255, 255, 0}
	LeaseRouter = []byte{10, 0, 0, 254}
	LeaseTime   = uint32(86400) // 1 Day
)

func Handle(s *dhcp.Server, message dhcp.Message) {
	log.Printf("DHCP message from NIC %s\n", message.ClientMAC)

	leaseTime := make([]byte, 4)
	binary.BigEndian.PutUint32(leaseTime, LeaseTime)

	var t byte

	if message.Type == dhcp.DHCPTypeDiscover {
		t = dhcp.DHCPTypeOffer
	} else if message.Type == dhcp.DHCPTypeRequest {
		t = dhcp.DHCPTypeACK
	} else {
		return
	}

	response := dhcp.NewMessage(t, message.TransactionID, ServerIP, LeaseIP, message.ClientMAC)
	response.SetOption(dhcp.OptionSubnetMask, LeaseMask)
	response.SetOption(dhcp.OptionRouter, LeaseRouter)
	response.SetOption(dhcp.OptionServerIdentifier, ServerIP)
	response.SetOption(dhcp.OptionIPAddressLeaseTime, leaseTime)

	s.BroadcastPacket(response.GetFrame())
}

func main() {
	fmt.Println("dhcpd (go-dhcp)")

	server, err := dhcp.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	server.HandleFunc(Handle)
	log.Fatal(server.ListenAndServe())
}
