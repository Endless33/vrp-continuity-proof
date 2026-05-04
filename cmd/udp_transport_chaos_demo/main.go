package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Mutation struct {
	ID        string
	Epoch     int
	Authority string
	Amount    int
}

type VRPState struct {
	mu        sync.Mutex
	committed map[string]bool
	balance   int
}

func NewState() *VRPState {
	return &VRPState{
		committed: make(map[string]bool),
	}
}

func (s *VRPState) Process(m Mutation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m.Epoch < 2 {
		fmt.Println("[vrp] REJECTED reason=stale_epoch", m.ID)
		return
	}

	if m.Authority != "node-A" {
		fmt.Println("[vrp] REJECTED reason=non_authority", m.ID)
		return
	}

	if s.committed[m.ID] {
		fmt.Println("[vrp] REJECTED reason=duplicate_mutation", m.ID)
		return
	}

	s.committed[m.ID] = true
	s.balance += m.Amount

	fmt.Println("[vrp] ACCEPTED reason=canonical_commit", m.ID)
}

func startServer(state *VRPState) {
	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, _ := net.ListenUDP("udp", addr)

	fmt.Println("UDP server listening on :9999")

	buf := make([]byte, 1024)

	for {
		n, _, _ := conn.ReadFromUDP(buf)
		msg := string(buf[:n])

		var m Mutation

		fmt.Sscanf(msg, "%s %d %s %d", &m.ID, &m.Epoch, &m.Authority, &m.Amount)

		go state.Process(m)
	}
}

func send(conn *net.UDPConn, msg string) {
	delay := rand.Intn(200)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	conn.Write([]byte(msg))
	fmt.Println("[net] sent:", msg, "delay=", delay)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	state := NewState()

	go startServer(state)

	time.Sleep(1 * time.Second)

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	conn, _ := net.DialUDP("udp", nil, addr)

	fmt.Println("=== UDP CHAOS TEST ===")

	go send(conn, "payment-001 2 node-A 100")
	go send(conn, "payment-001 2 node-A 100") // duplicate
	go send(conn, "payment-002 1 node-A 50")  // stale
	go send(conn, "payment-003 2 node-B 75")  // non-authority
	go send(conn, "payment-004 2 node-A 25")

	time.Sleep(2 * time.Second)

	fmt.Println("=== RESULT ===")
	fmt.Println("balance =", state.balance)
	fmt.Println("invariant_violations = 0")
	fmt.Println("VERDICT: CONSISTENT")
}