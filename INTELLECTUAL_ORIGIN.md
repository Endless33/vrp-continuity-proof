# VRP Intellectual Origin

Author: Vitalijus Riabovas

VRP (Veil Routing Protocol) and the continuity-first execution model were originally designed and developed by Vitalijus Riabovas.

This repository exists as a public, reproducible proof record for VRP execution correctness.

---

## Core Architectural Claims

VRP introduces and validates the following architectural principles:

- session identity is independent from transport
- transport failure must not imply session failure
- execution correctness is enforced at the commit layer
- authority is epoch-bound and monotonic
- stale authority cannot mutate state
- duplicate, replayed, or reordered inputs fail closed
- exactly one canonical mutation may commit under race
- endpoint changes must not reset logical session identity

---

## Verified Runtime Proofs

The following proof layers have been publicly demonstrated and reproduced.

### Oracle Proof Ladder

Validates clean, chaos, attack, and consensus behavior:

- Clean proof: session continuity under normal transition
- Chaos proof: duplicate / reordered / stale inputs rejected
- Attack proof: fake authority / replay / invalid epoch jump rejected
- Consensus proof: exactly one canonical winner under race

Proof ladder:

https://github.com/Endless33/vrp-continuity-proof/blob/main/oracle/ORACLE_10_PROOF_LADDER.md

Unified proof runner:

```bash
go run ./cmd/oracle_unified_proof_runner
Expected result:
Plain text
OVERALL VERDICT: CONTINUITY PRESERVED
UDP Continuity Handoff Proof
Validates that session identity is not bound to a UDP socket or endpoint.
Observed behavior:
UDP endpoint changed
session identity remained stable
replay was rejected
fake authority was rejected
execution continued without session reset
Proof:
https://github.com/Endless33/vrp-continuity-proof/blob/main/udp/UDP_CONTINUITY_HANDOFF_PROOF.md⁠�
Run server:
Bash
go run ./cmd/udp_continuity_server
Run client:
Bash
go run ./cmd/udp_continuity_client
Interpretation
VRP is not a reconnect mechanism.
VRP is an execution correctness layer for unreliable networks.
The central principle is:
Plain text
session != transport
A transport may fail, change, rebind, or be replaced.
The logical session must remain governed by authority, epoch, sequence, replay protection, and commit correctness.
Public Record
This document serves as a public authorship and architectural origin record for the VRP continuity-first execution model.
The architecture, invariants, execution model, and proof ladder represent original work and ongoing independent research by Vitalijus Riabovas.
Use of the code is governed by the LICENSE file.
Reproduction of the architectural design without attribution is discouraged.

## Commit message

```text
add VRP intellectual origin record

Markdown
## Intellectual Origin

The intellectual origin and architectural authorship record for VRP is documented here:

[INTELLECTUAL_ORIGIN.md](./INTELLECTUAL_ORIGIN.md)