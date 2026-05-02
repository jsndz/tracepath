package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

func main() {

	dns := os.Args[1]

	ips, err := net.LookupIP(dns)
	if err != nil {
		fmt.Println("couldn't find IP", err)
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)

	if err != nil {
		panic(err)
	}
	defer unix.Close(fd)

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("test"),
		},
	}
	packet, err := msg.Marshal(nil)

	if err != nil {
		panic(err)
	}
	ip := ips[0].To4()

	var b [4]byte
	copy(b[:], ip)
	// copy needs slice as arg
	// so convert b to slice using [:]
	//[:] is a slice operator
	// When applied to an array, it creates a slice that references the array’s underlying memory.
	start := time.Now()
	unix.Sendto(fd, packet, 0, &unix.SockaddrInet4{Addr: b})

	buf := make([]byte, 1500)

	for {
		n, from, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			panic(err)
		}

		ip := net.IP(from.(*unix.SockaddrInet4).Addr[:]).String()
		rtt := time.Since(start)

		ihl := int(buf[0]&0x0f) * 4
		if n < ihl {
			continue
		}

		msg, err := icmp.ParseMessage(1, buf[ihl:n])
		if err != nil {
			continue
		}

		switch body := msg.Body.(type) {
		case *icmp.Echo:
			fmt.Println("from", ip, "message:", string(body.Data), "rtt:", rtt)
		default:
			fmt.Println("from", ip, "type:", msg.Type, "rtt:", rtt)
		}
	}

}
