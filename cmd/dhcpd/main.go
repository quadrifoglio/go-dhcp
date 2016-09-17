package main

import (
	"fmt"
	"log"
	"net"

	"github.com/quadrifoglio/go-dhcp"
)

func HandleDiscover(s *dhcp.Server, transaction uint32, mac net.HardwareAddr) {
	log.Printf("DHCP Discover from NIC %s\n", mac)

	offer := dhcp.NewOffer(net.IPv4(10, 0, 0, 1), mac, transaction)
	offer.IP = []byte{10, 0, 0, 100}
	offer.Mask = []byte{255, 255, 255, 0}
	offer.Router = []byte{10, 0, 0, 1}

	s.BroadcastPacket(offer.GetBytes())
}

func HandleRequest(s *dhcp.Server, transaction uint32, mac net.HardwareAddr, requestedIP net.IP) {
	log.Printf("DHCP Request from NIC %s for IP address %s\n", mac, requestedIP)
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
