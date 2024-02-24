package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snaplen  = int32(320)
	promisc  = true
	timeout  = pcap.BlockForever
	filter   = "tcp and (tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18)"
	devFound = false
	results  = make(map[string]int)
)

func worker(host string, ports chan int) {
	for p := range ports {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, p), 1000*time.Millisecond)
		if err != nil {
			continue
		}
		conn.Close()

	}
}

func capture(iface, target string) {
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}

	defer handle.Close()
	filter := filter + "and src host " + target
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("Capturing packets")
	for packet := range source.Packets() {
		trnsportLayer := packet.TransportLayer()
		if trnsportLayer == nil {
			continue
		}
		srcPort := trnsportLayer.TransportFlow().Src().String()
		results[srcPort] += 1
	}
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln(" <capture_iface> <target_ip> <start port-end port>")
	}
	portsChan := make(chan int, 10)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	iface := os.Args[1]
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
		}
	}
	if !devFound {
		log.Panicf("Device named '%s' does not exist\n", iface)
	}
	ip := os.Args[2]
	ports := strings.Split(os.Args[3], "-")

	go capture(iface, ip)
	time.Sleep(1 * time.Second)

	startPort, startErr := strconv.Atoi(ports[0])
	if startErr != nil {
		fmt.Printf("start port is not a number")
		panic(startErr)
	}
	endPort, endErr := strconv.Atoi(ports[1])
	if endErr != nil {
		fmt.Printf("end port is not a number")
		panic(endErr)
	}
	for i := 0; i < cap(portsChan); i++ {
		go worker(ip, portsChan)
	}

	go func() {
		for i := startPort; i <= endPort; i++ {
			portsChan <- i
		}
	}()

	time.Sleep(2 * time.Second)
	for port, confidence := range results {
		if confidence >= 1 {
			fmt.Printf("Port %s open (confidence: %d)\n", port, confidence)
		}
	}
}
