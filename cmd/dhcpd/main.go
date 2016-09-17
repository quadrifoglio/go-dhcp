package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/quadrifoglio/go-dhcp"
)

func HandleDiscover(s *dhcp.Server, id uint32, mac net.HardwareAddr) {
	log.Printf("DHCP Discover from NIC %s\n", mac)

	serverIp := []byte{10, 0, 0, 1}
	clientIp := []byte{10, 0, 0, 100}
	clientMask := []byte{255, 255, 255, 0}
	clientRouter := []byte{10, 0, 0, 254}
	leaseTime := make([]byte, 4)

	binary.BigEndian.PutUint32(leaseTime, 86400)

	message := dhcp.NewMessage(dhcp.DHCPTypeOffer, id, serverIp, clientIp, mac)
	message.SetOption(dhcp.OptionSubnetMask, clientMask)
	message.SetOption(dhcp.OptionRouter, clientRouter)
	message.SetOption(dhcp.OptionServerIdentifier, serverIp)
	message.SetOption(dhcp.OptionIPAddressLeaseTime, leaseTime)

	s.BroadcastPacket(message.GetFrame())
}

func HandleRequest(s *dhcp.Server, id uint32, mac net.HardwareAddr, requestedIp net.IP) {
	log.Printf("DHCP Request from NIC %s for IP %s\n", mac, requestedIp)

	serverIp := []byte{10, 0, 0, 1}
	clientMask := []byte{255, 255, 255, 0}
	clientRouter := []byte{10, 0, 0, 254}
	leaseTime := make([]byte, 4)
	binary.BigEndian.PutUint32(leaseTime, 86400)

	message := dhcp.NewMessage(5, id, serverIp, requestedIp, mac)
	message.SetOption(dhcp.OptionSubnetMask, clientMask)
	message.SetOption(dhcp.OptionRouter, clientRouter)
	message.SetOption(dhcp.OptionServerIdentifier, serverIp)
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
