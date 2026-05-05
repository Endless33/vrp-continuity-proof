# VRP UDP Continuity Handoff Proof

Environment: Termux / Oracle-compatible runtime  
Transport: UDP  
Scenario: endpoint change (simulated path migration)

---

## PURPOSE

This demo validates that session identity is NOT bound to a UDP socket.

UDP endpoint may change.
Session identity MUST NOT.

---

## INVARIANTS

1. Session Identity Continuity  
   Session MUST survive endpoint change.

2. Replay Protection  
   Duplicate or old sequence MUST be rejected.

3. Authority Safety  
   Non-authoritative sender MUST be rejected.

4. Transport Independence  
   Session MUST remain valid across new UDP sockets.

---

## TEST SCENARIO

- client sends packets over UDP
- client rebinds to a new local port (simulated path change)
- server detects endpoint change
- session continues without reset
- replay and fake authority attempts are injected

---

## OBSERVED OUTPUT (REAL RUNTIME)

[NEW SESSION] session: session-xyz remote: 127.0.0.1:49570
[ACCEPTED] state: ACTIVE epoch: 1 seq: 1
[REJECTED] replay or old sequence
[PATH CHANGE DETECTED] old_remote: 127.0.0.1:49570 new_remote: 127.0.0.1:47546 session_identity_preserved=true
[ACCEPTED] state: MIGRATING seq: 2
[PATH CHANGE DETECTED] old_remote: 127.0.0.1:47546 new_remote: 127.0.0.1:50343 session_identity_preserved=true
[ACCEPTED] state: ACTIVE seq: 3
[REJECTED] fake or stale authority

---

## RESULT

- session_identity_preserved = true
- replay_rejected = true
- fake_authority_rejected = true
- endpoint_change_detected = true
- session_not_reset = true

---

## VERDICT

CONTINUITY PRESERVED

---

## INTERPRETATION

This is NOT a reconnect.

This is NOT session recreation.

This is:

session != transport

The UDP socket changed multiple times.  
The session did not.

---

## REPRODUCE

Run server:

go run ./cmd/udp_continuity_server

Run client:

go run ./cmd/udp_continuity_client

Works in:
- Oracle Linux
- Termux (Android)