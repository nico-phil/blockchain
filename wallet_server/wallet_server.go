package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Nico2220/blockchain/wallet"
)

type WalletServer struct {
	port int
	gateway string
}

func NewWalletServer(port int, gateway string) *WalletServer{
	return &WalletServer{port: port, gateway: gateway}
}

func(ws *WalletServer) Port() int{
	return ws.port
}

func(ws *WalletServer) GateWayString() string{
	return ws.gateway
}

func(ws *WalletServer) Index(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("error parsing file", err)
	}
	t.Execute(w, "")

}

type wrapper map[string]any

func(ws *WalletServer) CreateWallet(w http.ResponseWriter, r *http.Request){
	myWallet:= wallet.NewWallet()
	m, err:= myWallet.MarshalJSON()
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

func(ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request){
	var input wallet.TransactionRequest 
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("trasaction failed"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello stansactipn"))
}


func (ws *WalletServer) Run() error{
	router := http.NewServeMux()
	router.HandleFunc("/",  ws.Index)

	router.HandleFunc("POST /transactions", ws.CreateTransaction)
	router.HandleFunc("POST /wallet",ws.CreateWallet)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), router)
}
