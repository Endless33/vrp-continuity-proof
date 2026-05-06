# VRP Continuity Proof

Minimal executable proof of continuity-preserving execution under unstable transport conditions.

This repository demonstrates a core invariant:

A logical mutation may commit at most once per session.

The project evolved from a minimal commit-boundary proof into a live continuity-aware runtime capable of intercepting real OS packets through a Linux TUN interface and transporting them over a VRP-controlled UDP carrier.

---

# Core Principle

Traditional systems bind execution correctness to transport stability.

VRP separates them.

Session identity remains stable while transport paths may change.

Transport may fail.
Execution must not.

---

# What problem this solves

Real networks are unstable.

Packets are duplicated, reordered, delayed, replayed, rerouted, or dropped.

Traditional systems often respond with reconnect, renegotiation, session reset, duplicate execution, or reconciliation after the fact.

VRP introduces a different model:

- deterministic commit admission
- explicit authority ownership
- fail-closed validation
- transport-independent session continuity

The goal is not faster reconnect.

The goal is preserving execution correctness during transport instability.

---

# Repository Scope

This repository contains:

- continuity proof runtime
- deterministic commit model
- authority validation model
- replay-safe execution flow
- chaos execution scenarios
- transport continuity experiments
- real TUN + UDP runtime verification

This repository does NOT claim to be:

- a production VPN
- a finished distributed system
- a consensus replacement
- a complete cryptographic product

This is a continuity execution research runtime.

---

# Minimal Runtime Model

Core execution objects:

- session_id
- mutation_id
- authority
- epoch
- sequence_number
- commit_key

Commit rule:

commit_key = session_id + mutation_id

A mutation is accepted only if:

1. it originates from the current authority
2. its epoch is valid
3. its sequence is admissible
4. its commit key was never committed before

Everything else is rejected before state mutation.

---

# Fundamental Invariant

one logical mutation -> at most one commit

This invariant is enforced under:

- duplicate delivery
- retry after lost response
- stale authority
- stale epoch
- packet replay
- reordered arrival
- multi-node race candidates
- transport instability

---

# Quick Start

Run the base proof:

go run ./cmd/continuity_proof_demo

Expected:

VERDICT: CONSISTENT
Proof: no logical mutation committed more than once

---

# Demo Index

## Continuity Proof Demo

Run:

go run ./cmd/continuity_proof_demo

Demonstrates:

- duplicate rejection
- authority validation
- deterministic commit admission

Expected:

VERDICT: CONSISTENT

---

## Chaos Commit Demo

Run:

go run ./cmd/chaos_commit_demo

Simulates:

- duplicate delivery
- retry after lost response
- stale epoch
- non-authority mutation
- reordered input

Expected:

invariant_violations = 0
VERDICT: CONSISTENT

---

## Multi-Node Race Demo

Run:

go run ./cmd/multi_node_race_demo

Simulates a race condition between multiple nodes attempting to commit the same mutation.

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

candidate node=node-A -> ACCEPTED
candidate node=node-B -> REJECTED
committed_candidates=1
authority_conflicts=1
invariant_violations=0
VERDICT: CONSISTENT

Invariant:

one mutation -> at most one commit

This demonstrates deterministic convergence under concurrent execution.

---

## Chaos Orchestrator Demo

Run:

go run ./cmd/chaos_orchestrator_demo

Simulates combined failure conditions in a single execution:

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

committed_mutations=2
duplicates_rejected=2
stale_epoch_rejected=1
non_authority_rejected=1
race_candidates_rejected=1
invariant_violations=0
VERDICT: CONSISTENT

---

## UDP Transport Chaos Demo

Run:

go run ./cmd/udp_transport_chaos_demo

Uses a real UDP socket, goroutines, and random delivery delay.

Validates that unstable transport timing does not corrupt execution state.

Expected:

stale_epoch -> rejected
non_authority -> rejected
duplicate_mutation -> rejected
invariant_violations = 0
VERDICT: CONSISTENT

---

## Oracle Unified Proof Runner

Run:

go run ./cmd/oracle_unified_proof_runner

Expected:

CLEAN: CONSISTENT
CHAOS: CONSISTENT
ATTACK: CONSISTENT
CONSENSUS: CONSISTENT
OVERALL VERDICT: CONTINUITY PRESERVED

This runner validates:

- clean continuity correctness
- duplicate and reordered input rejection
- fake authority rejection
- replay rejection
- invalid epoch jump rejection
- deterministic canonical commit under race

---

## UDP Continuity Handoff Proof

Run server:

go run ./cmd/udp_continuity_server

Run client:

go run ./cmd/udp_continuity_client

This proof validates:

- real UDP endpoint changes
- session identity preservation
- replay rejection
- fake authority rejection
- continuity across socket changes

Expected server evidence:

PATH CHANGE DETECTED
session_identity_preserved=true
VERDICT: UDP CONTINUITY PRESERVED

---

# Run on Windows

Clone the repository and run:

run_vrp_proof_windows.bat

Or manually:

go run ./cmd/oracle_unified_proof_runner

For UDP continuity handoff proof, use two terminals.

Terminal 1:

go run ./cmd/udp_continuity_server

Terminal 2:

go run ./cmd/udp_continuity_client

---

# Stage 3A - Real TUN Integration

VRP moved beyond logical continuity proofs into real OS packet interception.

Verified on Oracle Linux 10 using:

- Linux TUN interface
- interface name: vrp0
- userspace Go runtime
- live ICMP traffic

Verified flow:

Linux kernel
-> TUN interface
-> VRP userspace runtime
-> packet inspection

Observed result:

ICMP packets entered the VRP runtime from the operating system network stack.

This verifies:

- real TUN creation
- kernel-to-userspace packet flow
- live packet interception
- runtime packet ownership

---

# Stage 3B - Real TUN + UDP Carrier Runtime

VRP successfully moved from local packet interception into real packet transport over a UDP carrier.

Verified on Oracle Linux 10 using:

- Linux TUN interface vrp0
- userspace VRP runtime
- UDP carrier transport
- live ICMP traffic
- real packet restoration back into the OS network stack

Verified runtime flow:

OS ping
-> TUN interface
-> VRP runtime
-> VRP frame encapsulation
-> UDP carrier transport
-> remote runtime processing
-> ICMP reply encapsulation
-> UDP return path
-> TUN write-back
-> live ping reply restored

Observed runtime evidence:

[LOCAL] VRP FRAME SENT OVER UDP
[REMOTE] VRP FRAME RECEIVED
[REMOTE] ICMP REPLY ENCAPSULATED
[REMOTE] UDP FRAME SENT BACK
[LOCAL] UDP FRAME RECEIVED
payload_written_to_tun=true
VERDICT: TUN PACKET RESTORED

Observed network result:

12 packets transmitted, 12 received, 0% packet loss

This verifies:

- real TUN packet interception
- VRP frame transport over UDP
- packet restoration into the OS network stack
- live continuity-capable userspace transport execution

Current status:

Stage 1 -> architectural model
Stage 2 -> continuity proof runtime
Stage 3A -> TUN integration
Stage 3B -> real UDP carrier transport verified

Next target:

Stage 3C - transport mutation and carrier handoff during live session execution.

Goal:

transport changes
session identity remains stable
execution continues

Transport may fail.
Execution must not.

---

# Stage 3C - Live UDP Carrier Handoff

VRP has successfully verified live carrier mutation during active packet execution.

A live ping session was running through:

OS ping
-> TUN interface
-> VRP runtime
-> VRP frame
-> UDP carrier-A
-> remote runtime
-> UDP return path
-> TUN write-back

During execution, the transport carrier was changed:

carrier-A:
127.0.0.1:12001

carrier-B:
127.0.0.1:13001

Observed handoff evidence:

[TRANSPORT FAILURE DETECTED]
old_carrier: carrier-A
old_remote: 127.0.0.1:12001
state: VOLATILE

[TRANSPORT REATTACH]
new_carrier: carrier-B
new_remote: 127.0.0.1:13001
session: session-xyz
session_identity_preserved=true
session_reset=false

After reattach, packet execution continued through carrier-B.

Observed network result:

16 packets transmitted, 16 received, 0% packet loss

This verifies:

- live TUN packet execution
- VRP frame transport over UDP
- carrier mutation during active session
- session identity preserved across carrier handoff
- no session reset
- continued ICMP execution after transport reattach

Current status:

Stage 1 -> architectural model
Stage 2 -> continuity proof runtime
Stage 3A -> TUN integration
Stage 3B -> real UDP carrier transport verified
Stage 3C -> live UDP carrier handoff verified

Next target:

Stage 3D - real external endpoint / network path mutation.

Goal:

real network changes
session identity remains stable
execution continues

Transport may fail.
Execution must not.

---

# Direction

Part of the VRP / Jumping VPN research:

- session identity above transport
- deterministic authority
- commit admission
- invariant-based runtime verification
- continuity under network uncertainty
- real packet execution through userspace runtime

---

# Core Idea

Continuity is not faster retry.

Continuity is a commit boundary.

Continuity is execution correctness preserved while transport changes.

---

# Intellectual Origin

VRP (Veil Routing Protocol) and the continuity-first execution model were originally designed and developed by Vitalijus Riabovas.

Core principles introduced:

- session != transport
- execution correctness over unreliable networks
- commit-layer authority model
- epoch-based authority transitions
- fail-closed mutation validation
- transport-independent session continuity

This repository provides reproducible runtime proofs of these concepts.

The architecture, invariants, and execution model are part of ongoing independent research.

Unauthorized reproduction of the design without attribution is discouraged.