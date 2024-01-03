package main

import (
	"bufio"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Unexpected error")

	}
	log.Printf("Recive %d bytes: %s\n", len(s), s)
	log.Println("Writing data")
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(s); err != nil {
		log.Fatalln("Unable to write data")
	}
	writer.Flush()
}

func main() {
	var port = "20080"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("Unable to bind port " + port)
	}
	log.Println("Listening on 0.0.0.0:" + port)
	for {
		conn, err := listener.Accept()
		log.Println("Recived connection")
		if err != nil {
			log.Fatalln("Unabele to accept connection")
		}
		go echo(conn)
	}
}
