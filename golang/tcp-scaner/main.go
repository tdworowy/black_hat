package main

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Dial("tcp", "scanme.npam.org:80")
	if err == nil {
		fmt.Printf("Connection successful")
	}
}
