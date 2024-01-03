package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func worker(host string, ports, results chan int) {
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, p))
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {

	args := os.Args[1:]
	if len(args) <= 1 {
		fmt.Printf("usage: <host> <start port>-<end port>")
		os.Exit(0)
	}
	portsChan := make(chan int, 300)
	resultsChan := make(chan int)

	host := args[0]
	ports := strings.Split(args[1], "-")
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
	var openPorts []int

	for i := 0; i < cap(portsChan); i++ {
		go worker(host, portsChan, resultsChan)
	}

	go func() {
		for i := startPort; i < endPort; i++ {
			portsChan <- i
		}
	}()

	for i := startPort; i < endPort; i++ {
		port := <-resultsChan
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(portsChan)
	close(resultsChan)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}
