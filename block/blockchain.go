package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Nico2220/blockchain/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER      = 20
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

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timeStamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
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
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              int
	mu                sync.Mutex
}

func NewBlockchain(blockchainAddress string, port int) *Blockchain {
	b := &Block{}
	bc := &Blockchain{}
	bc.CreateBlock(0, b.Hash())
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	return bc
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarsalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
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

func (bc *Blockchain) CreateTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	// Todo

	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	// if bc.VerifyTransaction(senderPublicKey, s, t) {
	// 	if bc.CalculateTotalAmount(sender) < value{
	// 		log.Println("ERROR:", "Not wnougth balance in wallet")
	// 		return false
	// 	}
	// 	bc.transactionPool = append(bc.transactionPool, t)
	// 	return true
	// }

	bc.transactionPool = append(bc.transactionPool, t)

	return true

	// log.Println("ERROR:", "Verify transaction")

	// return false
}

func (bc *Blockchain) VerifyTransaction(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	v := ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
	return v
}

func (bc *Blockchain) Copytransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value))

	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{timeStamp: 0, nonce: nonce, previousHash: previousHash, transactions: transactions}
	guessBlockStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessBlockStr[:difficulty] == zeros

}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.Copytransactions()
	previousHash := bc.LasBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.transactionPool) == 0 {
		return false
	}

	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LasBlock().previousHash
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=MINING", "status=success")
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, c := range bc.chain {
		for _, t := range c.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}

	return totalAmount
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(senderBlockchainAddress, recipientBlockchainAdress string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress,
		recipientBlockchainAdress,
		value,
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blochain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})

}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address:       %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipient_blockchain_address:     %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value: %.1f\n", t.value)
}

type TransactionRequest struct {
	SenderBlockchainAddress  *string  `json:"sender_blochain_address"`
	RecipientBlochainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey          *string  `json:"sender_public_key"`
	Value                    *float32 `json:"value"`
	Signature                *string  `json:"signature"`
}

func (t *TransactionRequest) Validate() map[string]string {
	errorText := "missing value"
	mapError := map[string]string{}
	if t.SenderBlockchainAddress == nil {
		mapError["sender_blockchain_address"] = errorText
	}

	if t.RecipientBlochainAddress == nil {
		mapError["recipient_blockchain_address"] = errorText
	}

	if t.SenderPublicKey == nil {
		mapError["sender_public_key"] = errorText
	}

	if t.Value == nil {
		mapError["value"] = errorText
	}

	if t.Signature == nil {
		mapError["signature"] = errorText
	}
	return mapError
}
