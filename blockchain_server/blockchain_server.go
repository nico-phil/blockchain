package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Nico2220/blockchain/block"
	"github.com/Nico2220/blockchain/wallet"
)

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
		log.Printf("private_key %v", minerWallet.PrivateKeyStr())
		log.Printf("public_key %v", minerWallet.PublicKeyStr())
		log.Printf("private_key %v", minerWallet.BlockchainAddress())
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

func (bsc *BlockchainServer) Run() error {
	router := http.NewServeMux()
	// router.HandleFunc("/", HelloWorld)

	router.HandleFunc("/chain", bsc.GetChainHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", bsc.port), router)
}
