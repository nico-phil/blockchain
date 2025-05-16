package main

import (
	"os"
	"strconv"
)

func main(){
	port, _ := strconv.Atoi(os.Getenv("port"))

	app := NewBlockchainServer(port)
	
	app.Run()
	
}