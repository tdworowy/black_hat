package main

import (
	"fmt"
	"net"
	"sync"
)

func worker(ports chan int, wg *sync.WaitGroup) {
	for p := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("scanme.npam.org:%d", p))
		if err != nil {
			continue
		}
		conn.Close()
		fmt.Println(fmt.Sprintf("Connection successfulon on:%d", p))
		wg.Done()
	}
}

// func worker(ports chan int, wg *sync.WaitGroup) {
// 	for p := range ports {
// 		fmt.Println(p)
// 		wg.Done()
// 	}
// }

func main() {
	ports := make(chan int, 100)
	var wg sync.WaitGroup
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg)
	}
	for i := 1; i < 1024; i++ {
		wg.Add(1)
		ports <- i
	}
	wg.Wait()
	close(ports)
}
