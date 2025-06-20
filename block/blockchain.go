package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	BLOCKCHAIN_PORT_RANGE_START        = 5001
	BLOCKCHAIN_PORT_RANGE_END          = 5003
	NEIGHBOR_IP_RANGE_START            = 0
	NEIGHTBOR_IP_RANGE_END             = 1
	BLOCKCHAIN_NEIGTHBOR_SYNC_TIME_SEC = 20
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

func(b *Block) Nonce() int {
	return b.nonce
}

func(b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func(b *Block) Transactions() []*Transaction {
	return b.transactions
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

func(b *Block) UnMarshalJSON(data []byte) error{
	var previousHash string
	v := struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timeStamp,
		Nonce:        b.nonce,
		PreviousHash: previousHash,
		Transactions: b.transactions,
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ph, _ := hex.DecodeString(v.PreviousHash)
	copy(b.previousHash[:], ph[:3])

	return nil
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

	neighbors    []string
	muxNeighbors sync.Mutex
}

func NewBlockchain(blockchainAddress string, port int) *Blockchain {
	b := &Block{}
	bc := &Blockchain{}
	bc.CreateBlock(0, b.Hash())
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	return bc
}

func(b *Blockchain) Chain() []*Block {
	return b.chain
}

func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(utils.GetHost(), bc.port,
		NEIGHBOR_IP_RANGE_START, NEIGHTBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START,
		BLOCKCHAIN_PORT_RANGE_END)

	log.Printf("neighbors:%v", bc.neighbors)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()

	time.AfterFunc(BLOCKCHAIN_NEIGTHBOR_SYNC_TIME_SEC*time.Second, bc.StartSyncNeighbors)
}
func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
	bc.ResolveConfilcts()
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarsalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

func(bc *Blockchain) UnMarsalJSON(data []byte) error {
	v := struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
		resp, _ := client.Do(req)
		fmt.Printf("resp:%v", resp)
	}
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
	client := http.Client{}
	if isTransacted {
		for _, n := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			tr := TransactionRequest{
				SenderBlockchainAddress:  &sender,
				RecipientBlochainAddress: &recipient,
				SenderPublicKey:          &publicKeyStr,
				Value:                    &value,
				Signature:                &signatureStr,
			}
			m, _ := json.Marshal(tr)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			buf := bytes.NewBuffer(m)
			
			req, _ := http.NewRequest(http.MethodPut, endpoint, buf)
			resp, _ := client.Do(req)
			fmt.Printf("resp:%v", resp)

		}
	}

	return isTransacted
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	// if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
	// 	// if bc.CalculateTotalAmount(sender) < value{
	// 	// 	log.Println("ERROR:", "Not wnougth balance in wallet")
	// 	// 	return false
	// 	// }
	// 	bc.transactionPool = append(bc.transactionPool, t)
	// 	return true
	// }

	bc.transactionPool = append(bc.transactionPool, t)

	return true

	// log.Println("ERROR:", "Verify transaction")

	// return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, endpoint, nil)
		response, _ := client.Do(req)
		fmt.Printf("%v", response)

	}
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

func(bc *Blockchain) ValidChain(chain []*Block) bool{
	prevBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		currentBlock := chain[currentIndex]
		if currentBlock.PreviousHash() != prevBlock.Hash() {
			return false
		}

		// if !bc.ValidProof(currentBlock.Nonce(), currentBlock.PreviousHash(), currentBlock.Transactions(), MINING_DIFFICULTY){
		// 	return false
		// }

		prevBlock = currentBlock
		currentIndex++
	}

	return true
}

func(bc *Blockchain) ResolveConfilcts() bool{
	var longuestChain []*Block  = nil
	maxLength := len(longuestChain)
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		response, _ := http.Get(endpoint)
		if response.StatusCode == http.StatusOK {
			var bcResp Blockchain
			_ = json.NewDecoder(response.Body).Decode(&bcResp)
			
			chain := bcResp.Chain()
			//&& bc.ValidChain(chain)
			if len(chain) > maxLength  {
				longuestChain = chain
				maxLength = len(chain)
			}
		}
	}

	if longuestChain != nil {
		bc.chain = longuestChain
		fmt.Println("conflic resolve, longest chain:", longuestChain)
		return true
	}

	return false
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

func(t *Transaction) UnMarshalJSON(data []byte) error {
	v := struct {
		Sender    *string  `json:"sender_blochain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &t.senderBlockchainAddress,
		Recipient: &t.recipientBlockchainAddress,
		Value:     &t.value,
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return nil
	}

	return nil
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

type AmountResponse struct {
	Amount float32 `json:"amount"`
}
