package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Node      string
	SessionID string
	MutationID string
	Epoch     uint64
	Amount    int
}

type RaceRuntime struct {
	SessionID            string
	CurrentEpoch         uint64
	CurrentAuthority     string
	CommittedMutations   map[string]bool
	Balance              int
	CommittedCandidates  int
	RejectedCandidates   int
	AuthorityConflicts   int
	InvariantViolations  int
	Winner               string
}

func NewRaceRuntime() *RaceRuntime {
	return &RaceRuntime{
		SessionID:          "session-abc",
		CurrentEpoch:       51,
		CurrentAuthority:   "node-A",
		CommittedMutations: make(map[string]bool),
	}
}

func (r *RaceRuntime) ResolveRace(candidates []Candidate) {
	fmt.Println("=== VRP MULTI-NODE RACE DEMO ===")
	fmt.Println("Invariant: for a given (session, epoch, mutation), only one candidate may commit")
	fmt.Println()

	fmt.Println("RACE INPUT:")
	for _, c := range candidates {
		fmt.Printf("  candidate node=%s session=%s epoch=%d mutation=%s amount=%d\n",
			c.Node, c.SessionID, c.Epoch, c.MutationID, c.Amount)
	}
	fmt.Println()

	fmt.Println("STEP 1 -> sort candidates deterministically")
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Node < candidates[j].Node
	})
	for _, c := range candidates {
		fmt.Printf("  candidate_order=%s\n", c.Node)
	}
	fmt.Println()

	fmt.Println("STEP 2 -> commit gate evaluation")
	for _, c := range candidates {
		r.ProcessCandidate(c)
	}
}

func (r *RaceRuntime) ProcessCandidate(c Candidate) {
	fmt.Printf("[vrp] evaluating node=%s mutation=%s epoch=%d\n", c.Node, c.MutationID, c.Epoch)

	if c.SessionID != r.SessionID {
		r.RejectedCandidates++
		fmt.Println("[vrp] decision=REJECTED reason=session_mismatch")
		fmt.Println()
		return
	}

	if c.Epoch != r.CurrentEpoch {
		r.RejectedCandidates++
		fmt.Println("[vrp] decision=REJECTED reason=stale_or_invalid_epoch")
		fmt.Println()
		return
	}

	if c.Node != r.CurrentAuthority {
		r.RejectedCandidates++
		r.AuthorityConflicts++
		fmt.Println("[vrp] decision=REJECTED reason=non_canonical_authority")
		fmt.Println()
		return
	}

	if r.CommittedMutations[c.MutationID] {
		r.RejectedCandidates++
		r.AuthorityConflicts++
		fmt.Println("[vrp] decision=REJECTED reason=already_committed")
		fmt.Println()
		return
	}

	r.CommittedMutations[c.MutationID] = true
	r.CommittedCandidates++
	r.Balance += c.Amount
	r.Winner = c.Node

	fmt.Println("[vrp] decision=ACCEPTED reason=canonical_commit")
	fmt.Println()
}

func (r *RaceRuntime) Validate(expectedBalance int) {
	if r.CommittedCandidates != 1 {
		r.InvariantViolations++
	}

	if r.Balance != expectedBalance {
		r.InvariantViolations++
	}

	if r.Winner != r.CurrentAuthority {
		r.InvariantViolations++
	}
}

func main() {
	runtime := NewRaceRuntime()

	candidates := []Candidate{
		{
			Node:       "node-B",
			SessionID:  "session-abc",
			MutationID: "transfer-777",
			Epoch:     51,
			Amount:    100,
		},
		{
			Node:       "node-A",
			SessionID:  "session-abc",
			MutationID: "transfer-777",
			Epoch:     51,
			Amount:    100,
		},
	}

	runtime.ResolveRace(candidates)
	runtime.Validate(100)

	fmt.Println("=== PROOF REPORT ===")
	fmt.Printf("session=%s\n", runtime.SessionID)
	fmt.Printf("epoch=%d\n", runtime.CurrentEpoch)
	fmt.Printf("canonical_authority=%s\n", runtime.CurrentAuthority)
	fmt.Printf("winner=%s\n", runtime.Winner)
	fmt.Printf("committed_candidates=%d\n", runtime.CommittedCandidates)
	fmt.Printf("rejected_candidates=%d\n", runtime.RejectedCandidates)
	fmt.Printf("authority_conflicts=%d\n", runtime.AuthorityConflicts)
	fmt.Printf("balance=%d\n", runtime.Balance)
	fmt.Printf("invariant_violations=%d\n", runtime.InvariantViolations)
	fmt.Println()

	if runtime.InvariantViolations == 0 {
		fmt.Println("VERDICT: CONSISTENT")
		fmt.Println("Proof: multi-node race converged to one canonical commit")
	} else {
		fmt.Println("VERDICT: BROKEN")
		fmt.Println("Proof failed: multiple candidates affected state")
	}
}