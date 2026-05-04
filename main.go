package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

func main() {
	var (
		maxHops = flag.Int("h", 30, "max hops")
		timeout = flag.Int("t", 2, "timeout in seconds")
		probes  = flag.Int("p", 3, "probes per hop")
	)
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("usage: traceroute [options] <host>")
		flag.PrintDefaults()
		return
	}

	dns := flag.Arg(0)

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
	MAXHOPS := *maxHops
	TIMEOUT := *timeout
	PROBES := *probes
	for i := 0; i < MAXHOPS; i++ {
		unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_TTL, i+1)
		unix.SetsockoptTimeval(fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO,
			&unix.Timeval{Sec: int32(TIMEOUT), Usec: 0})
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
		done := false
		fmt.Printf("%2d  ", i+1)
		for try := 0; try < PROBES; try++ {
			start := time.Now()
			unix.Sendto(fd, packet, 0, &unix.SockaddrInet4{Addr: b})

			buf := make([]byte, 1500)

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
			reply, err := icmp.ParseMessage(1, buf[ihl:n])
			if err != nil {
				fmt.Printf("* ")
				continue
			}

			fmt.Printf("%s %v ", ip, rtt)

			if reply.Type == ipv4.ICMPTypeEchoReply {
				done = true
				break
			}
		}
		fmt.Println()

		if done {
			break
		}

	}
}
