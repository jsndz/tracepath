# 1: get Ip from DNS

use net.LookupIP
// Your program asks OS resolver.
// OS queries configured DNS server.
// DNS returns one or more IPs.
// Go gives them as []net.IP.

```go
	dns := os.Args[1]

	ips, err := net.LookupIP(dns)
	if err != nil {
		fmt.Println("couldn't find IP", err)
	}

```

# 2: sending ICMP packet

ICMP is a control and diagnostic protocol used with IP.
used for:
Report network errors
Indicate delivery problems
Test reachability
Measure network path behavior

OPen a raw socket
create icmp message
send to the ip
here ip is in net.Ip convert it to 4 bytes
then copy the bytes to array

```go
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
	unix.Sendto(fd, packet, 0, &unix.SockaddrInet4{Addr: b})

```


# 3: configuring ttl, timeout, timer


so the main thing is like we send a echo message through icmp protocol,
and it travels through the router and reaches the destination 
for example, google server
send the packet it will hop and finally reach the server
using raw socket to send the msg


we need to get how many hops/ routers will the message go through to reach server 
so use ttl,
set it up using: 		
```go
unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_TTL, ttl)
```

since the code is synchronus, it will wait for message continously 
so lets add time out

```go
unix.SetsockoptTimeval(fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO,
			&unix.Timeval{Sec: 2})
```

also added latency for checking how much time each reply takes

code:
```go

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


```