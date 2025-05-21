package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Nico2220/blockchain/block"
	"github.com/Nico2220/blockchain/utils"
	"github.com/Nico2220/blockchain/wallet"
)

type WalletServer struct {
	port    int
	gateway string
}

func NewWalletServer(port int, gateway string) *WalletServer {
	return &WalletServer{port: port, gateway: gateway}
}

func (ws *WalletServer) Port() int {
	return ws.port
}

func (ws *WalletServer) GateWayString() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("error parsing file", err)
	}
	t.Execute(w, "")

}

type wrapper map[string]any

func (ws *WalletServer) CreateWallet(w http.ResponseWriter, r *http.Request) {
	myWallet := wallet.NewWallet()
	m, err := myWallet.MarshalJSON()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		e := wrapper{"error": "error marshaling"}
		me, _ := json.Marshal(e)
		w.Write(me)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input wallet.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": err.Error()})
		return
	}

	// validate transaction
	if !input.ValidTransaction() {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, wrapper{"error": "cannot process transaction entity"})
		return
	}

	publicKey := utils.PublickKeyFromString(*input.SenderPublickKey)
	privateKey := utils.PrivateKeyFromString(*input.SenderPrivateKey, publicKey)
	value, err := strconv.ParseFloat(*input.Value, 32)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, wrapper{"error": err.Error()})
	}
	value32 := float32(value)

	transaction := wallet.NewTransaction(privateKey, publicKey, *input.SenderBlockchainAddress, *input.RecepientBlockchainAddress, value32)
	signature := transaction.GenerateSignature()
	signatureStr := signature.String()

	bt := block.TransactionRequest{
		SenderBlockchainAddress:  input.SenderBlockchainAddress,
		RecipientBlochainAddress: input.RecepientBlockchainAddress,
		SenderPublicKey:          input.SenderPublickKey,
		Value:                    &value32,
		Signature:                &signatureStr,
	}

	fmt.Println("block_transaction:", bt)

	utils.WriteJSON(w, http.StatusOK, wrapper{"transaction": input})
}

func (ws *WalletServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("/", ws.Index)

	router.HandleFunc("POST /transactions", ws.CreateTransaction)
	router.HandleFunc("POST /wallet", ws.CreateWallet)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), router)
}
