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