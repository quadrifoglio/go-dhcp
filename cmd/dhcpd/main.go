package main

import (
	"fmt"
	"log"
	"net"

	"github.com/quadrifoglio/go-dhcp"
)

func HandleDiscover(socket net.PacketConn, transaction uint32, mac net.HardwareAddr) {
	log.Printf("DHCP Discover from NIC %s\n", mac)

	offer := dhcp.NewOffer(net.IPv4(10, 0, 0, 1), mac, transaction)
	offer.IP = net.IPv4(10, 0, 0, 100)
	offer.Mask = net.IPv4(255, 255, 255, 0)
	offer.Router = net.IPv4(10, 0, 0, 1)

	ip := &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 68}

	_, err := socket.WriteTo(offer.GetBytes(), ip)
	if err != nil {
		log.Println(err)
	}
}

func HandleRequest(socket net.PacketConn, transaction uint32, mac net.HardwareAddr, requestedIP net.IP) {
	log.Printf("DHCP Request from NIC 0x%x for IP address %s\n", mac, requestedIP)
}

func HandleRelease(socket net.PacketConn, mac net.HardwareAddr) {
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

	log.Fatal(server.Start())
}
