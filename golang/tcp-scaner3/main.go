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
	ports_chan := make(chan int, 300)
	results_chan := make(chan int)

	host := args[0]
	ports := strings.Split(args[1], "-")
	start_port, start_err := strconv.Atoi(ports[0])
	if start_err != nil {
		fmt.Printf("start port is not a number")
		panic(start_err)
	}
	end_port, end_err := strconv.Atoi(ports[1])
	if end_err != nil {
		fmt.Printf("end port is not a number")
		panic(end_err)
	}
	var open_ports []int

	for i := 0; i < cap(ports_chan); i++ {
		go worker(host, ports_chan, results_chan)
	}

	go func() {
		for i := start_port; i < end_port; i++ {
			ports_chan <- i
		}
	}()

	for i := start_port; i < end_port; i++ {
		port := <-results_chan
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}

	close(ports_chan)
	close(results_chan)
	sort.Ints(open_ports)
	for _, port := range open_ports {
		fmt.Printf("%d open\n", port)
	}
}
