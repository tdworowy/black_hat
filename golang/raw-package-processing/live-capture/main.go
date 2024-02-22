package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	iface    = "eth0"
	snaplen  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	filter   = "tcp and port 80"
	devFound = false
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device naped %s does not exist\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}
	defer handle.Close()
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		hex_data := hex.Dump(packet.Data())
		fmt.Print(packet)
		fmt.Println(hex_data)
	}
}
