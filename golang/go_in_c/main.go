package main

import "C"
import "fmt"

//export Start
func Start() {
	fmt.Println("Yo from Go")
}

func main() {}

// go build -buildmode=c-archive
