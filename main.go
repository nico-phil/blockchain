package main

import (
	// "github.com/Nico2220/blockchain/block"
	"fmt"

	"github.com/Nico2220/blockchain/wallet"
)

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
}
