package main

import (
	"black_hat_go/metasploit/rpc"
	"fmt"
	"log"
)

func main() {
	host := "172.25.97.90:55552" //os.Getenv("MSFHOST")
	pass := "test"               //os.Getenv("MSFPASS")
	user := "msf"

	if host == "" || pass == "" {
		log.Fatalln("Missing required environment variable MSFHOST or MSFPASS")
	}
	msf, err := rpc.New(host, user, pass)
	if err != nil {
		log.Panicln(err)
	}
	defer msf.Logout()

	sessions, err := msf.SessionList()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Sessions:")
	for _, session := range sessions {
		fmt.Printf("%5d  %s\n", session.ID, session.Info)
	}
}
