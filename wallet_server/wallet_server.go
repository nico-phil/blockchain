package main

import (
	"fmt"
	"html/template"
	"net/http"
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


func (ws *WalletServer) Run() error{
	router := http.NewServeMux()
	// router.HandleFunc("/", HelloWorld)

	router.HandleFunc("/wallet", ws.Index)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), router)
}
