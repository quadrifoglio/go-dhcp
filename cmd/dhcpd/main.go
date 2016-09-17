package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/quadrifoglio/go-dhcp"
)

var (
	ServerIP = []byte{10, 0, 0, 1}

	LeaseIP     = []byte{10, 0, 0, 100}
	LeaseMask   = []byte{255, 255, 255, 0}
	LeaseRouter = []byte{10, 0, 0, 254}
	LeaseTime   = uint32(86400) // 1 Day
)

func HandleDiscover(s *dhcp.Server, id uint32, mac net.HardwareAddr) {
	log.Printf("DHCP Discover from NIC %s\n", mac)

	leaseTime := make([]byte, 4)
	binary.BigEndian.PutUint32(leaseTime, LeaseTime)

	message := dhcp.NewMessage(dhcp.DHCPTypeOffer, id, ServerIP, LeaseIP, mac)
	message.SetOption(dhcp.OptionSubnetMask, LeaseMask)
	message.SetOption(dhcp.OptionRouter, LeaseRouter)
	message.SetOption(dhcp.OptionServerIdentifier, ServerIP)
	message.SetOption(dhcp.OptionIPAddressLeaseTime, leaseTime)

	s.BroadcastPacket(message.GetFrame())
}

func HandleRequest(s *dhcp.Server, id uint32, mac net.HardwareAddr, requestedIp net.IP) {
	log.Printf("DHCP Request from NIC %s for IP %s\n", mac, requestedIp)

	leaseTime := make([]byte, 4)
	binary.BigEndian.PutUint32(leaseTime, LeaseTime)

	message := dhcp.NewMessage(dhcp.DHCPTypeACK, id, ServerIP, requestedIp, mac)
	message.SetOption(dhcp.OptionSubnetMask, LeaseMask)
	message.SetOption(dhcp.OptionRouter, LeaseRouter)
	message.SetOption(dhcp.OptionServerIdentifier, ServerIP)
	message.SetOption(dhcp.OptionIPAddressLeaseTime, leaseTime)

	s.BroadcastPacket(message.GetFrame())
}

func HandleRelease(s *dhcp.Server, mac net.HardwareAddr) {
	log.Printf("DHCP Release from NIC 0x%x\n", mac)
}

func main() {
	fmt.Println("dhcpd (go-dhcp)")

	server, err := dhcp.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	server.HandleDiscover(HandleDiscover)
	server.HandleRequest(HandleRequest)
	server.HandleRelease(HandleRelease)

	log.Fatal(server.ListenAndServe())
}
