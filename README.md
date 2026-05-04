# Traceroute Clone

## Overview

A command-line utility built in Go that traces the network path from your machine to a destination host.

It works by sending probe packets with increasing TTL values and collecting ICMP responses from routers along the route.

This reveals every hop between source and destination.

---

## Features

- Trace route to hostname or IP
- Per-hop latency measurement
- Timeout handling
- Configurable max hops
- Multiple probes per hop
- Detect final destination reached
- Clean CLI output
- IPv4 support

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

## Tech Stack

- Go

---

## Run

```bash id="d0x1ea"
sudo env "PATH=$PATH" go run main.go -h 20 -t 1 -p 5 google.com
```
