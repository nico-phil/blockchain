package main

import (
	"fmt"

	"github.com/Nico2220/blockchain/utils"
)

func main(){
	isFound := utils.IsFoundHost("127.0.0.1", 5002)
	fmt.Println(isFound)
}