package main

import (
	"fmt"

	"github.com/Endless33/vrp-continuity-proof/internal/proof"
)

func main() {
	fmt.Println("=== VRP CONTINUITY PROOF DEMO ===")
	fmt.Println("Invariant: one logical mutation may commit at most once per session")
	fmt.Println()

	events := []proof.Event{
		{
			Name:       "original payment",
			SessionID:  "sess-42",
			MutationID: "pay-001",
			NodeID:     "node-A",
			Authority:  "node-A",
			Action:     "commit",
			Amount:     100,
		},
		{
			Name:       "retry after lost response",
			SessionID:  "sess-42",
			MutationID: "pay-001",
			NodeID:     "node-A",
			Authority:  "node-A",
			Action:     "commit",
			Amount:     100,
		},
		{
			Name:       "same mutation from non-authority",
			SessionID:  "sess-42",
			MutationID: "pay-001",
			NodeID:     "node-B",
			Authority:  "node-A",
			Action:     "commit",
			Amount:     100,
		},
		{
			Name:       "new valid payment",
			SessionID:  "sess-42",
			MutationID: "pay-002",
			NodeID:     "node-A",
			Authority:  "node-A",
			Action:     "commit",
			Amount:     50,
		},
	}

	engine := proof.NewEngine()

	for index, event := range events {
		fmt.Printf("EVENT %d → %s\n", index+1, event.Name)

		beforeTraceCount := len(engine.Trace)
		engine.Apply(event)

		for _, line := range engine.Trace[beforeTraceCount:] {
			fmt.Println("  →", line)
		}

		fmt.Println()
	}

	engine.Verify()
	proof.PrintReport(engine)
}