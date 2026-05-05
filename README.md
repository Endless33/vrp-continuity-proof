# VRP Continuity Proof

Minimal executable proof of a continuity-aware commit boundary.

This repository demonstrates one core invariant:

> A logical mutation may commit at most once per session.

---

## Problem

In distributed systems, a network failure can turn one logical action into multiple executions.

Client sends a payment → server commits → response is lost → client retries.

The same logical mutation appears twice.

Without a commit boundary, it may execute twice.

Result:

- duplicate payments  
- duplicate orders  
- inventory drift  
- conflicting state  
- reconciliation after the fact  

---

## Model

This proof uses a minimal runtime model:

- `session_id`
- `mutation_id`
- `authority`
- `commit_key = session_id + mutation_id`

A mutation is accepted only if:

1. it comes from the current authority  
2. its commit key has not been committed before  

Everything else is rejected before it mutates state.

---

## Run

```bash
go run ./cmd/continuity_proof_demo
```

---

## Expected result

```
VERDICT: CONSISTENT
Proof: no logical mutation committed more than once
```

---

## What this proves

This proof demonstrates commit admission under:

- retry after lost response  
- duplicate logical mutation  
- non-authoritative proposal  

Invariant:

```
one logical mutation → at most one commit
```

---

## What this does NOT claim

- This is not a production VPN  
- This is not a full distributed consensus system  
- This is not Raft / Spanner replacement  

This is a minimal proof artifact.

---

## Chaos Commit Demo

Run:

```bash
go run ./cmd/chaos_commit_demo

This demo simulates:
duplicate delivery
retry after lost response
stale epoch
non-authority mutation
reordered input
Expected:
one logical mutation commits at most once
invalid inputs are rejected
invariant_violations = 0
Example output:

VERDICT: CONSISTENT

---

## Multi-Node Race Demo

Run:

```bash
go run ./cmd/multi_node_race_demo
```

This demo simulates a race condition between multiple nodes attempting to commit the same mutation.

Scenario:

- same session  
- same epoch  
- same mutation  
- different nodes  

Expected behavior:

- multiple candidates may exist  
- only one candidate is allowed to commit  
- all others are rejected before state mutation  

Example outcome:

```
candidate node=node-A → ACCEPTED
candidate node=node-B → REJECTED

committed_candidates=1
authority_conflicts=1
invariant_violations=0

VERDICT: CONSISTENT
```

Invariant:

```
one mutation → at most one commit
```

This demonstrates deterministic convergence under concurrent execution.

---

## Chaos Orchestrator Demo

Run:

```bash
go run ./cmd/chaos_orchestrator_demo
```

This demo simulates combined failure conditions in a single execution:

- duplicate delivery  
- delayed retry  
- reordered arrival  
- stale epoch  
- multi-node race candidate  
- non-authority mutation  

Expected behavior:

- valid mutations commit once  
- duplicate inputs are rejected  
- stale and invalid inputs are rejected  
- concurrent candidates do not create divergence  

Example outcome:

```
committed_mutations=2
duplicates_rejected=2
stale_epoch_rejected=1
non_authority_rejected=1
race_candidates_rejected=1
invariant_violations=0

VERDICT: CONSISTENT
```

Invariant:

```
one mutation → at most one commit
```

This demonstrates system-level correctness under combined network chaos.

---

## UDP Transport Chaos Demo

Run:

```bash
go run ./cmd/udp_transport_chaos_demo

This demo uses a real UDP socket, goroutines, and random delivery delay.
It validates that unstable transport timing does not corrupt execution state.
Expected:

stale_epoch → rejected
non_authority → rejected
duplicate_mutation → rejected
invariant_violations = 0
VERDICT: CONSISTENT

---

## Run on Windows

Clone the repository and run:

```bat
run_vrp_proof_windows.bat

Or manually:

Bash
go run ./cmd/oracle_unified_proof_runner

For UDP continuity handoff proof, use two terminals:

Terminal 1:

Bash
go run ./cmd/udp_continuity_server

Terminal 2:

Bash
go run ./cmd/udp_continuity_client

---

## Direction

Part of the VRP / Jumping VPN research:

- session identity above transport  
- deterministic authority  
- commit admission  
- invariant-based runtime verification  
- continuity under network uncertainty  

---

## Core idea

Continuity is not faster retry.

Continuity is a commit boundary.

---

## Intellectual Origin

VRP (Veil Routing Protocol) and the continuity-first execution model
were originally designed and developed by Vitalijus Riabovas.

Core principles introduced:

- session != transport
- execution correctness over unreliable networks
- commit-layer authority model
- epoch-based authority transitions
- fail-closed mutation validation

This repository provides reproducible runtime proofs of these concepts.

The architecture, invariants, and execution model are part of ongoing independent research.

Unauthorized reproduction of the design without attribution is discouraged.