package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		fmt.Println(device.Name)
		fmt.Printf("Description: %s\n", device.Description)
		fmt.Printf("Flags: %d\n", device.Flags)
		for _, address := range device.Addresses {
			fmt.Printf("  IP: %s\n", address.IP)
			fmt.Printf("  Netmask: %s\n", address.Netmask)

		}
	}
}
