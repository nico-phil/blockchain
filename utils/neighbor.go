package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func IsFoundHost(host string, port int) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(target)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		log.Printf("%s %v\n", target, err)
		return false
	}

	return true
}
func FindNeighbors(myHost string, myPort, startIp, endIp, startPort, endPort int) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for port := startPort; port <= endPort; port += 1 {
		for ip := startIp; ip <= endIp; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}


func GetHost() string {
	defaultHost := "127.0.0.1"
	hostname, err := os.Hostname()
	if err != nil {
		return defaultHost
	}
	addresses, err := net.LookupHost(hostname)
	if err != nil {
		return defaultHost
	}
	return addresses[0]
}