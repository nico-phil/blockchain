package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
)

type Block struct {
	timeStamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timeStamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timeStamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func (b *Block) Print() {
	fmt.Printf("timeStamp    %d\n", b.timeStamp)
	fmt.Printf("previousHash    %x\n", b.previousHash)
	fmt.Printf("nonce    %d\n", b.nonce)
	for _, transaction := range b.transactions {
		transaction.Print()
	}
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := &Blockchain{}
	bc.CreateBlock(0, b.hash())
	return bc

}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LasBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s chain %d %s\n ", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}

	fmt.Printf("%s\n", strings.Repeat("=", 25))
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) Copytransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value))

	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHas [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{timeStamp: 0, nonce: nonce, previousHash: previousHas, transactions: transactions}
	guessBlockStr := fmt.Sprintf("%x", guessBlock.hash())
	return guessBlockStr[:difficulty] == zeros

}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.Copytransactions()
	previousHas := bc.LasBlock().hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHas, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(senderBlockchainAddress, recipientBlockchainAdress string, value float32) *Transaction {
	return &Transaction{senderBlockchainAddress, recipientBlockchainAdress, value}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string
		Recipient string
		value     float32
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		value:     t.value,
	})

}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address:       %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient_blockchain_address:     %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value: %.1f\n", t.value)
}

func main() {
	bc := NewBlockchain()
	// bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	previousHash := bc.LasBlock().hash()
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, previousHash)
	// bc.Print()

	bc.AddTransaction("A", "B", 2.1)
	bc.AddTransaction("x", "y", 1.1)
	previousHash = bc.LasBlock().hash()
	nonce = bc.ProofOfWork()
	bc.CreateBlock(nonce, previousHash)
	bc.Print()
}
