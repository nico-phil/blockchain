package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Nico2220/blockchain/block"
	"github.com/Nico2220/blockchain/utils"
	"github.com/Nico2220/blockchain/wallet"
)

type Wrapper map[string]any

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port int
}

func NewBlockchainServer(port int) *BlockchainServer {
	return &BlockchainServer{port: port}
}

func (bcs *BlockchainServer) Port() int {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minerWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("miner_wallet_private_key %v", minerWallet.PrivateKeyStr())
		log.Printf("miner_wallet_public_key %v", minerWallet.PublicKeyStr())
		log.Printf("miner_blockchain_address %v", minerWallet.BlockchainAddress())
	}

	return bc
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func (bcs *BlockchainServer) GetChainHandler(w http.ResponseWriter, r *http.Request) {
	bc := bcs.GetBlockchain()
	js, err := bc.MarsalJSON()
	if err != nil {
		log.Fatal(err)
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (bcs *BlockchainServer) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("transaction from blochain server")
	var t block.TransactionRequest
	err := utils.ReadJSON(r, &t)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"error": err.Error()})
		return
	}

	if mapError := t.Validate(); len(mapError) > 0 {
		fmt.Println("not validate")
		utils.WriteJSON(w, http.StatusUnprocessableEntity, Wrapper{"error": mapError})
		return
	}

	fmt.Println("tr_after_validate", *t.RecipientBlochainAddress)

	publicKey := utils.PublickKeyFromString(*t.SenderPublicKey)
	signature := utils.SignatureFromString(*t.Signature)

	bc := bcs.GetBlockchain()

	isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlochainAddress,
		*t.Value, publicKey, signature)
	fmt.Println("iscreated:", isCreated)

	if !isCreated {
		utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"transaction": "transaction is not created"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, Wrapper{"transaction": t})
}

func (bcs *BlockchainServer) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	bc := bcs.GetBlockchain()
	transactions := bc.TransactionPool()

	var t struct {
		Transactions []*block.Transaction `json:"transactions"`
		Length       int                  `json:"length"`
	}

	t.Transactions = transactions
	t.Length = len(transactions)

	utils.WriteJSON(w, http.StatusOK, Wrapper{"transaction": t.Transactions, "length": t.Length})
}

func (bcs *BlockchainServer) Mine(w http.ResponseWriter, r *http.Request) {
	bc := bcs.GetBlockchain()
	isMined := bc.Mining()

	if !isMined {
		utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"error": "mine failed"})
		return
	}

	utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"message": "mine succed"})
}

func (bcs *BlockchainServer) StartMining(w http.ResponseWriter, r *http.Request) {
	bc := bcs.GetBlockchain()
	bc.StartMining()

	utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"message": "start mining"})
}

func(bcs *BlockchainServer) GetAmount(w http.ResponseWriter, r *http.Request) {
	bc := bcs.GetBlockchain()
	blockchainAddress := r.URL.Query().Get("blockchain_address")

	if blockchainAddress == "" {
		utils.WriteJSON(w, http.StatusNotFound, Wrapper{"error": "missing blockchain address"})
		return
	}
	amount := bc.CalculateTotalAmount(blockchainAddress)

	utils.WriteJSON(w, http.StatusOK, Wrapper{"amount": amount})
}

func (bsc *BlockchainServer) Run() error {
	fmt.Println("blockchain_server running on:", bsc.port)
	router := http.NewServeMux()
	// router.HandleFunc("/", HelloWorld)

	router.HandleFunc("/transactions", bsc.GetTransactionHandler)
	router.HandleFunc("POST /transactions", bsc.TransactionHandler)
	router.HandleFunc("/chain", bsc.GetChainHandler)
	router.HandleFunc("/mine", bsc.Mine)
	router.HandleFunc("/mine/start", bsc.StartMining)
	router.HandleFunc("/amount", bsc.GetAmount)
	return http.ListenAndServe(fmt.Sprintf(":%d", bsc.port), router)
}
