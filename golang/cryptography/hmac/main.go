package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var key = []byte("some random key")

func checkMAC(message, revMAC []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	calcMac := mac.Sum(nil)

	return hmac.Equal(calcMac, revMAC)
}

func main() {
	message := []byte("The red eagle flies at 10:00")
	mac, _ := hex.DecodeString("69d2c7b6fbbfcaeb72a3172f4662601d1f16acfb46339639ac8c10c8da64631d")

	if checkMAC(message, mac) {
		fmt.Println("EQUAL")
	} else {
		fmt.Println("NOT EQUAL")
	}

}
