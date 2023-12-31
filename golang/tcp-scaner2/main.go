package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("scanme.npam.org:%d", p))
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var open_ports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i < 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(open_ports)
	for _, port := range open_ports {
		fmt.Printf("%d open\n", port)
	}
}
