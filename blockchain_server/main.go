package main

import (
	"log"
	"os"
	"strconv"
)

func init() {

}

func main() {
	port, _ := strconv.Atoi(os.Getenv("port"))

	app := NewBlockchainServer(port)

	err := app.Run()
	if err != nil {
		log.Fatal("error starting server", err)
	}

}
