# Traceroute Clone

## Overview

A command-line utility built in Go that traces the network path from your machine to a destination host.

It works by sending probe packets with increasing TTL values and collecting ICMP responses from routers along the route.

This reveals every hop between source and destination.

---

## Why This Project

This project helps understand the Network Layer deeply through real packet behavior.

Topics covered:

* IP routing
* TTL / hop limits
* Router forwarding
* ICMP error messages
* Round-trip latency
* DNS resolution
* Raw sockets
* Real-world internet paths

---

## Features

* Trace route to hostname or IP
* Per-hop latency measurement
* Timeout handling
* Configurable max hops
* Multiple probes per hop
* Detect final destination reached
* Clean CLI output
* IPv4 support

---

## How It Works

1. Resolve domain to destination IP
2. Send first probe with TTL = 1
3. First router decrements TTL to zero
4. Router sends ICMP Time Exceeded reply
5. Record router IP + response time
6. Send next probe with TTL = 2
7. Continue until destination replies or max hops reached

---

## Example Output

```bash id="y5hj9k"
$ traceroute google.com

1   192.168.1.1        2ms
2   10.0.0.1           7ms
3   172.16.12.4       14ms
4   142.250.x.x       26ms
```

---

## Tech Stack

* Go
* `net`
* `golang.org/x/net/icmp`
* `golang.org/x/net/ipv4`
* Goroutines
* CLI

---

## Project Structure

```bash id="b4cv6n"
cmd/
internal/
  ├── traceroute/
  ├── probe/
  ├── listener/
  ├── dns/
  ├── output/
  └── timer/
```

---

## Challenges

* Raw socket permissions may require sudo/admin
* Some routers block ICMP
* Firewalls may hide hops
* Network paths can vary dynamically

---

## Future Improvements

* IPv6 support
* Reverse DNS lookup
* Geo-location of hops
* Packet loss statistics
* Parallel probing
* JSON output mode
* TUI visualization

---

## Run

```bash id="d0x1ea"
sudo go run main.go google.com
```

