package utils

import (
	"fmt"
)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(target)

	// _, err := net.DialTimeout()

	return false
}