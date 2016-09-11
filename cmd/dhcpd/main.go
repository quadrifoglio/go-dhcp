package main

import (
	"fmt"
	"log"

	"github.com/quadrifoglio/go-dhcp"
)

func HandleDiscover(info dhcp.Info) {
	log.Printf("DHCP Discover from NIC 0x%x\n", info.MAC)
}

func HandleRequest(info dhcp.Info) {
	log.Printf("DHCP Request from NIC 0x%x for IP address %s\n", info.MAC, info.RequestIP)
}

func HandleRelease(info dhcp.Info) {
	log.Printf("DHCP Release from NIC 0x%x\n", info.MAC)
}

func main() {
	fmt.Println("dhcpd (go-dhcp)")

	server, err := dhcp.NewServer("10.0.0.0", 10)
	if err != nil {
		log.Fatal(err)
	}

	server.HandleFunc("discover", HandleDiscover)
	server.HandleFunc("request", HandleRequest)
	server.HandleFunc("release", HandleRelease)

	log.Fatal(server.Start())
}
