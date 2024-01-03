package main

import (
	"io"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 512)

	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		log.Printf("Recive %d bytes: %s\n", size, string(b))
		log.Println("Writing data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
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
