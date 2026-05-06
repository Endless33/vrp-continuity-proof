# VRP Stage Progress

This document tracks the execution progress of the VRP (Veil Routing Protocol) continuity runtime.

The goal is not to claim a production-ready VPN.

The goal is to verify continuity invariants step-by-step under increasingly realistic network conditions.

---

# Core Invariant

Transport may fail.

Execution must not.

---

# Stage 1 — Architectural Model

Status: VERIFIED

Focus:

- protocol architecture
- session identity model
- transport abstraction
- continuity semantics
- state machine
- authority model

Verified concepts:

- session != transport
- deterministic transitions
- explicit recovery states
- transport replacement model

Artifacts:

- architectural specification
- protocol documentation
- state machine definition
- invariants definition

Result:

Architectural continuity model established.

---

# Stage 2 — Continuity Proof Runtime

Status: VERIFIED

Focus:

- commit correctness
- authority validation
- replay rejection
- duplicate mutation prevention
- deterministic commit admission

Verified runtime invariants:

- one mutation commits at most once
- stale authority rejected
- duplicate mutation rejected
- invalid mutation rejected
- deterministic convergence preserved

Verified scenarios:

- retry after lost response
- duplicate delivery
- stale epoch
- non-authority mutation
- multi-node race
- UDP transport disorder

Representative verdict:

VERDICT: CONSISTENT

Result:

Continuity-aware commit boundary verified.

---

# Stage 3A — TUN Integration

Status: VERIFIED

Focus:

- Linux TUN interface
- OS packet interception
- userspace runtime injection

Verified behavior:

- TUN interface created
- packets intercepted from OS stack
- runtime packet visibility confirmed

Verified environment:

- Oracle Linux 10
- Go userspace runtime
- Linux TUN driver

Result:

VRP connected to real OS packet flow.

---

# Stage 3B — UDP Carrier Runtime

Status: VERIFIED

Focus:

- VRP frame encapsulation
- UDP carrier transport
- packet restoration into TUN

Verified runtime flow:

OS ping
→ TUN
→ VRP frame
→ UDP transport
→ remote runtime
→ ICMP encapsulation
→ UDP return path
→ TUN write-back

Observed runtime evidence:

- VRP frame transport verified
- ICMP reply restoration verified
- payload_written_to_tun=true

Observed network result:

12 packets transmitted
12 received
0% packet loss

Result:

Real userspace transport runtime verified.

---

# Stage 3C — Carrier Handoff

Status: VERIFIED

Focus:

- transport carrier replacement
- continuity during carrier mutation

Verified behavior:

carrier-A
→ carrier-B
→ same session identity
→ execution continued

Observed runtime evidence:

[TRANSPORT FAILURE DETECTED]

[TRANSPORT REATTACH]

session_identity_preserved=true
session_reset=false

Result:

Live carrier handoff verified without session reset.

---

# Stage 3D — UDP Endpoint Rebinding

Status: VERIFIED

Focus:

- live UDP endpoint mutation
- endpoint rebinding during active execution

Verified behavior:

old endpoint
→ new endpoint
→ same session
→ ICMP execution continued

Observed runtime evidence:

[REMOTE ENDPOINT MUTATION DETECTED]

old_remote:
127.0.0.1:56356

new_remote:
127.0.0.1:56477

session_identity_preserved=true
session_reset=false

Observed network result:

12 packets transmitted
12 received
0% packet loss

Result:

Live endpoint rebinding verified during active runtime execution.

---

# Current Runtime State

Current verified capabilities:

- session continuity model
- deterministic commit runtime
- authority validation
- replay rejection
- TUN packet interception
- UDP userspace transport
- carrier handoff
- endpoint rebinding
- ICMP continuity restoration
- runtime evidence generation

Current runtime properties:

- observable behavior
- explicit state transitions
- fail-closed mutation handling
- deterministic runtime decisions
- continuity-aware execution

---

# Not Yet Verified

The following remain future targets:

- real external network path mutation
- Wi-Fi ↔ LTE handoff
- NAT rebinding across real networks
- encrypted production transport
- distributed authority coordination
- multi-node runtime replication
- production-grade tunnel runtime
- kernel integration
- QUIC transport runtime
- formal verification

---

# Next Target — Stage 3E

Goal:

Real external network mutation.

Target scenario:

Wi-Fi
→ hotspot / LTE
→ public endpoint mutation
→ session survives
→ execution continues

Focus:

- real NAT rebinding
- real route mutation
- real path instability
- continuity across physical network changes

Expected invariant:

transport changed
session identity preserved
execution continues

---

# Philosophy

VRP does not treat transport instability as exceptional.

Transport volatility is treated as a normal runtime condition.

The runtime adapts while preserving session continuity and execution correctness.

---

# Summary

VRP currently demonstrates:

- continuity-aware runtime execution
- transport-independent session identity
- deterministic mutation admission
- observable continuity behavior
- live userspace transport mutation handling

This repository tracks the transition from architectural continuity research into observable runtime behavior.

Transport may fail.

Execution must not.