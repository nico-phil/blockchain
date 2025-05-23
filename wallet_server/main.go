package main

import (
	"log"
	"os"
	"strconv"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("port"))

	walletServer := NewWalletServer(port, "http://localhost:5001")

	err := walletServer.Run()
	if err != nil {
		log.Fatal("error starting the server", err)
	}
}
