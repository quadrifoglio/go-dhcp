package main

import (
	"fmt"
	"log"

	"github.com/quadrifoglio/go-dhcp"
)

func main() {
	fmt.Println("dhcpd (go-dhcp)")

	server, err := dhcp.NewServer("10.0.0.0", 10, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}
