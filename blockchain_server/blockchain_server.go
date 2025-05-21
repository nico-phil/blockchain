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

func(bcs *BlockchainServer) TransactionHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("transactio from blochain server")
	var t block.TransactionRequest
	err := utils.ReadJSON(r, &t)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, Wrapper{"error": err.Error()})
		return 
	}

	fmt.Println(t)

	utils.WriteJSON(w, http.StatusCreated, Wrapper{"transaction": t})
}

func (bsc *BlockchainServer) Run() error {
	fmt.Println("blockchain_server running on:", bsc.port)
	router := http.NewServeMux()
	// router.HandleFunc("/", HelloWorld)

	router.HandleFunc("POST /transactions", bsc.TransactionHandler)
	router.HandleFunc("/chain", bsc.GetChainHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", bsc.port), router)
}
