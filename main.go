package main

import (
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce int
	previousHash string
	timeStamp int64
	transactions []string
}

func NewBlock(nonce int, previousHash string ) *Block {
	return &Block{
		nonce: nonce,
		previousHash: previousHash,
		timeStamp: time.Now().UnixNano(),
	}
}

func(b *Block) Print(){
	fmt.Printf("timeStamp    %d\n", b.timeStamp)
	fmt.Printf("previousHash    %s\n", b.previousHash)
	fmt.Printf("nonce    %d\n", b.nonce)
	fmt.Printf("transactions    %v\n", b.transactions)
}

type Blockchain struct {
	transactionPool []string
	chain []*Block
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{}
	bc.CreateBlock(0, "Init hash")
	return bc
  
}


func(bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func(bc *Blockchain) Print(){
	for i, block := range bc.chain {
		fmt.Printf("%s chain %d %s\n ", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}

	fmt.Printf("%s\n", strings.Repeat("=", 25))
}

func main(){
	bc := NewBlockchain()
	bc.CreateBlock(5, "has 1")
	bc.CreateBlock(10, "has 2")
	bc.Print()
}