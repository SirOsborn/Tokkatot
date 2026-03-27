package main

import (
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	// Generate VAPID keys
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("================================================================")
	fmt.Println("TOKKATOT VAPID KEY GENERATOR")
	fmt.Println("================================================================")
	fmt.Println("Copy these into your .env file for Push Notifications:")
	fmt.Println("")
	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", publicKey)
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", privateKey)
	fmt.Println("")
	fmt.Println("Keep these safe! Never share your private key.")
	fmt.Println("================================================================")
}
