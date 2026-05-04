package main

import (
	"fmt"
	"sort"
)

type Mutation struct {
	ID        string
	Session  string
	Epoch    uint64
	Node     string
	Amount   int
	Label    string
	DelayMS  int
	Sequence int
}

type Runtime struct {
	Session              string
	CurrentEpoch         uint64
	Authority            string
	Committed            map[string]bool
	Balance              int
	CommittedCount        int
	DuplicatesRejected    int
	StaleRejected         int
	NonAuthorityRejected  int
	DelayedHandled        int
	ReorderedHandled      int
	RaceRejected          int
	InvariantViolations   int
}

func NewRuntime() *Runtime {
	return &Runtime{
		Session:       "session-chaos-001",
		CurrentEpoch:  7,
		Authority:     "node-A",
		Committed:     make(map[string]bool),
	}
}

func (r *Runtime) Process(m Mutation) {
	fmt.Printf("[delivery] seq=%d delay=%dms label=%s\n", m.Sequence, m.DelayMS, m.Label)
	fmt.Printf("[vrp] session=%s epoch=%d node=%s mutation=%s amount=%d\n",
		m.Session, m.Epoch, m.Node, m.ID, m.Amount)

	if m.DelayMS > 0 {
		r.DelayedHandled++
		fmt.Println("[vrp] note=delayed_input_revalidated")
	}

	if m.Session != r.Session {
		r.InvariantViolations++
		fmt.Println("[vrp] decision=REJECTED reason=session_mismatch")
		fmt.Println()
		return
	}

	if m.Epoch < r.CurrentEpoch {
		r.StaleRejected++
		fmt.Println("[vrp] decision=REJECTED reason=stale_epoch")
		fmt.Println()
		return
	}

	if m.Node != r.Authority {
		r.NonAuthorityRejected++
		r.RaceRejected++
		fmt.Println("[vrp] decision=REJECTED reason=non_canonical_authority")
		fmt.Println()
		return
	}

	if r.Committed[m.ID] {
		r.DuplicatesRejected++
		fmt.Println("[vrp] decision=REJECTED reason=duplicate_mutation")
		fmt.Println()
		return
	}

	r.Committed[m.ID] = true
	r.Balance += m.Amount
	r.CommittedCount++

	fmt.Println("[vrp] decision=ACCEPTED reason=canonical_commit")
	fmt.Println()
}

func (r *Runtime) Validate() {
	if r.CommittedCount != 2 {
		r.InvariantViolations++
	}

	if r.Balance != 125 {
		r.InvariantViolations++
	}

	if r.Committed["payment-001"] != true {
		r.InvariantViolations++
	}

	if r.Committed["payment-004"] != true {
		r.InvariantViolations++
	}
}

func main() {
	fmt.Println("=== VRP CHAOS ORCHESTRATOR DEMO ===")
	fmt.Println("Invariant: unstable delivery must not corrupt execution state")
	fmt.Println()

	fmt.Println("CHAOS CONDITIONS:")
	fmt.Println("- duplicate delivery")
	fmt.Println("- delayed retry")
	fmt.Println("- reordered arrival")
	fmt.Println("- stale epoch")
	fmt.Println("- multi-node race candidate")
	fmt.Println("- non-authority rejection")
	fmt.Println()

	rt := NewRuntime()

	events := []Mutation{
		{
			ID:        "payment-001",
			Session:   "session-chaos-001",
			Epoch:     7,
			Node:      "node-A",
			Amount:    100,
			Label:     "original mutation",
			DelayMS:   0,
			Sequence:  1,
		},
		{
			ID:        "payment-004",
			Session:   "session-chaos-001",
			Epoch:     7,
			Node:      "node-A",
			Amount:    25,
			Label:     "valid mutation delivered early due to reorder",
			DelayMS:   20,
			Sequence:  4,
		},
		{
			ID:        "payment-001",
			Session:   "session-chaos-001",
			Epoch:     7,
			Node:      "node-A",
			Amount:    100,
			Label:     "retry after lost response",
			DelayMS:   120,
			Sequence:  2,
		},
		{
			ID:        "payment-002",
			Session:   "session-chaos-001",
			Epoch:     6,
			Node:      "node-A",
			Amount:    50,
			Label:     "stale epoch arrives late",
			DelayMS:   200,
			Sequence:  3,
		},
		{
			ID:        "payment-003",
			Session:   "session-chaos-001",
			Epoch:     7,
			Node:      "node-B",
			Amount:    75,
			Label:     "multi-node race candidate from non-authority",
			DelayMS:   40,
			Sequence:  5,
		},
		{
			ID:        "payment-004",
			Session:   "session-chaos-001",
			Epoch:     7,
			Node:      "node-A",
			Amount:    25,
			Label:     "duplicate of reordered mutation",
			DelayMS:   260,
			Sequence:  6,
		},
	}

	fmt.Println("NETWORK DELIVERY ORDER:")
	for _, e := range events {
		fmt.Printf("  delivered seq=%d label=%s\n", e.Sequence, e.Label)
	}
	fmt.Println()

	fmt.Println("LOGICAL ORDER CHECK:")
	sorted := make([]Mutation, len(events))
	copy(sorted, events)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Sequence < sorted[j].Sequence
	})

	for i, e := range events {
		if e.Sequence != i+1 {
			rt.ReorderedHandled++
			fmt.Printf("  reorder_detected delivered_position=%d logical_seq=%d mutation=%s\n",
				i+1, e.Sequence, e.ID)
		}
	}
	fmt.Println()

	fmt.Println("RUNTIME PROCESSING:")
	for _, e := range events {
		rt.Process(e)
	}

	rt.Validate()

	fmt.Println("=== PROOF REPORT ===")
	fmt.Printf("session=%s\n", rt.Session)
	fmt.Printf("epoch=%d\n", rt.CurrentEpoch)
	fmt.Printf("authority=%s\n", rt.Authority)
	fmt.Printf("committed_mutations=%d\n", rt.CommittedCount)
	fmt.Printf("duplicates_rejected=%d\n", rt.DuplicatesRejected)
	fmt.Printf("stale_epoch_rejected=%d\n", rt.StaleRejected)
	fmt.Printf("non_authority_rejected=%d\n", rt.NonAuthorityRejected)
	fmt.Printf("race_candidates_rejected=%d\n", rt.RaceRejected)
	fmt.Printf("delayed_inputs_handled=%d\n", rt.DelayedHandled)
	fmt.Printf("reordered_inputs_detected=%d\n", rt.ReorderedHandled)
	fmt.Printf("balance=%d\n", rt.Balance)
	fmt.Printf("invariant_violations=%d\n", rt.InvariantViolations)
	fmt.Println()

	if rt.InvariantViolations == 0 {
		fmt.Println("VERDICT: CONSISTENT")
		fmt.Println("Proof: duplicate, delay, reorder, stale epoch, and race candidate did not corrupt state")
	} else {
		fmt.Println("VERDICT: BROKEN")
		fmt.Println("Proof failed: chaos conditions corrupted state")
	}
}