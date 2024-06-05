package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce int
	previousHash [32]byte
	timeStamp int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte ) *Block {
	return &Block{
		nonce: nonce,
		previousHash: previousHash,
		timeStamp: time.Now().UnixNano(),
	}
}

func(b *Block) hash() [32]byte{
	m, _ := json.Marshal(b)
	fmt.Print(m)
	return sha256.Sum256(m)
}

func(b* Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp int64 `json:"timestamp"`
		Nonce int `json:"nonce"`
		PreviousHash [32]byte `json:"previousHash"`
		Transaction []string  `json:"transaction"`
	}{
		Timestamp : b.timeStamp,
		Nonce: b.nonce,
		PreviousHash: b.previousHash,
		Transaction: b.transactions,
	})
}

func(b *Block) Print(){
	fmt.Printf("timeStamp    %d\n", b.timeStamp)
	fmt.Printf("previousHash    %x\n", b.previousHash)
	fmt.Printf("nonce    %d\n", b.nonce)
	fmt.Printf("transactions    %v\n", b.transactions)
}

type Blockchain struct {
	transactionPool []string
	chain []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := &Blockchain{}
	bc.CreateBlock(0, b.hash())
	return bc
  
}


func(bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
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
	 NewBlockchain()
	
	// bc.Print()
}