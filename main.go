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

	ip := ips[0].To4()

	var b [4]byte
	copy(b[:], ip)
	// copy needs slice as arg
	// so convert b to slice using [:]
	//[:] is a slice operator
	// When applied to an array, it creates a slice that references the array’s underlying memory.
	for i := 0; i < 30; i++ {
		unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_TTL, i+1)
		unix.SetsockoptTimeval(fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO,
			&unix.Timeval{Sec: 2})
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i,
				Data: []byte("test"),
			},
		}
		packet, err := msg.Marshal(nil)

		if err != nil {
			panic(err)
		}
		start := time.Now()
		unix.Sendto(fd, packet, 0, &unix.SockaddrInet4{Addr: b})

		buf := make([]byte, 1500)

		n, from, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			panic(err)
		}

		ip := net.IP(from.(*unix.SockaddrInet4).Addr[:]).String()
		rtt := time.Since(start)
		// header length of 4 bytes, so multiply by 4 to get the actual header length in bytes
		ihl := int(buf[0]&0x0f) * 4
		//if total lenght is less than header length invalid
		if n < ihl {
			continue
		}
		// parse icmp messsage
		reply, err := icmp.ParseMessage(1, buf[ihl:n])
		if err != nil {
			continue
		}
		// if echo message
		switch body := reply.Body.(type) {
		case *icmp.Echo:
			fmt.Println("from", ip, "message:", string(body.Data), "rtt:", rtt, "hop:", i+1)
		default:
			fmt.Println("from", ip, "type:", reply.Type, "rtt:", rtt, "hop:", i+1)
		}
		// if you get echo reply,  reached the destination and can stop
		if reply.Type == ipv4.ICMPTypeEchoReply {
			break
		}
	}
}
