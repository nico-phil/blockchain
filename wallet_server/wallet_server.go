package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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

func (ws *WalletServer) GateWay() string {
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

	m, _ := json.Marshal(bt)
	buf := bytes.NewBuffer(m)

	fmt.Println("gateway:", ws.GateWay())
	rep, err := http.Post(ws.GateWay()+"/transactions", "application/json", buf)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": err.Error()})
	}
	if rep.StatusCode != http.StatusCreated {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "transaction failed :("})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, wrapper{"message": "success"})

}

func (ws *WalletServer) GetAmount(w http.ResponseWriter, r *http.Request) {

	client := http.Client{}
	endpoint := fmt.Sprintf("%s/amount", ws.gateway)

	blockchainAddress := r.URL.Query().Get("blockchain_address")

	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	q := req.URL.Query()
	q.Add("blockchain_address", blockchainAddress)
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "cannot get amount"})
		return
	}

	if response.StatusCode != http.StatusOK {
		utils.WriteJSON(w, response.StatusCode, wrapper{"error": "cannot get amount"})
		return
	}

	var a block.AmountResponse
	err = json.NewDecoder(response.Body).Decode(&a)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, wrapper{"error": "cannot decode amount reponse"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, wrapper{"amount": a.Amount})

}

func (ws *WalletServer) Run() error {
	fmt.Println("wallet_server running on:", ws.port)
	router := http.NewServeMux()
	router.HandleFunc("/", ws.Index)

	router.HandleFunc("POST /transactions", ws.CreateTransaction)
	router.HandleFunc("POST /wallet", ws.CreateWallet)
	router.HandleFunc("GET /wallet/amount", ws.GetAmount)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), router)
}
