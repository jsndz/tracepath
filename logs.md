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
