# 1: get Ip from DNS

use net.LookupIP
// Your program asks OS resolver.
// OS queries configured DNS server.
// DNS returns one or more IPs.
// Go gives them as []net.IP.
s