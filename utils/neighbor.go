package utils

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

var PATTERN = regexp.MustCompile()

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(target)

	_, err := net.DialTimeout("tcp", target, 1 * time.Second)
	if err != nil {
		log.Printf("%s %v\n", target, err)
		return false
	}

	return true
}

func FindNeigbors(myHost string, myPort, startIp, endIp, startPort, endPort int) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	m := PATTERN.FindStringSubmatch(myHost)
	if m != nil {
		return nil
	}

	return []string{}
}