package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	dns := os.Args[1]

	ips, err := net.LookupIP(dns)
	if err != nil {
		fmt.Println("couldn't find IP", err)
	}

	for _, ip := range ips {
		fmt.Println(ip)
	}

}
