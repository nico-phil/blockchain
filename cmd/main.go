package main

import (
	"fmt"

	"github.com/Nico2220/blockchain/utils"
)

func main() {
	isFound := utils.FindNeighbors("127.0.0.1", 5001, 0, 3, 5001, 5003)
	fmt.Println(isFound)
}
