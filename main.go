package main

import (
	"github.com/Nico2220/blockchain/block"
)

func main() {
	// w := wallet.NewWallet()
	// fmt.Println(w.PrivateKeyStr())
	// fmt.Println(w.PublicKeyStr())

	// fmt.Println(w.BlockchainAddress())

	// t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	// fmt.Printf("signature %s\n", t.GenerateSignature())

	// myBlockchainAddress := "my_blockchain_address"
	// blockchain := block.NewBlockchain(myBlockchainAddress)


	// blockchain.AddTransaction("A", "B", 1.0)
	// blockchain.Mining()
	

	// blockchain.AddTransaction("C", "D", 2.0)
	// blockchain.AddTransaction("X", "X", 3.0)
	// blockchain.Mining()

	// blockchain.Print()

	// fmt.Printf("my %.1f\n", blockchain.CalculateTotalAmount(myBlockchainAddress))
	// fmt.Printf("C %.1f\n", blockchain.CalculateTotalAmount("C"))
	// fmt.Printf("D %.1f\n", blockchain.CalculateTotalAmount("D"))

	
	blockChain := block.NewBlockchain("my_address")

	blockChain.AddTransaction("A", "B", 1.0, nil, nil)
	blockChain.Mining()

	

	blockChain.Print()



}
