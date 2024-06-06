package main

import (
	"github.com/Nico2220/blockchain/block"
)

func main() {
	address := "MY_BLOCKCHAIN_ADDRESS"
	bc := block.NewBlockchain(address)
	// bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	bc.Mining()

	bc.AddTransaction("A", "B", 2.1)
	bc.AddTransaction("x", "y", 1.1)
	bc.Mining()
	bc.Print()

	// fmt.Printf("%.1f\n", bc.calculateTotalAmount(address))
	// fmt.Printf("A%.1f\n", bc.calculateTotalAmount("A"))
	// fmt.Printf("C%.1f\n", bc.calculateTotalAmount("C"))
	// fmt.Printf("y%.1f\n", bc.calculateTotalAmount("y"))
	// fmt.Printf("x%.1f\n", bc.calculateTotalAmount("x"))
}