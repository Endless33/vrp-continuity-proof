=== VRP ORACLE CONSENSUS PROOF DEMO ===
Runtime: 2026-05-05T11:21:17Z
Invariant: one logical mutation may commit at most once
Invariant: only current authority may commit
Invariant: loser proposals must not mutate state

[DELIVERY ORDER]
delay=1ms id=proposal-stale mutation=payment-commit-001 epoch=0 authority=node-a priority=100
delay=5ms id=proposal-fake mutation=payment-commit-001 epoch=1 authority=node-x priority=99
delay=10ms id=proposal-b mutation=payment-commit-001 epoch=1 authority=node-a priority=20
delay=30ms id=proposal-a mutation=payment-commit-001 epoch=1 authority=node-a priority=10

[PROPOSAL] proposal-stale
REJECTED: stale or invalid epoch
accepted=0 rejected=1 committed=map[]

[PROPOSAL] proposal-fake
REJECTED: non-authoritative proposal
accepted=0 rejected=2 committed=map[]

[PROPOSAL] proposal-b
ACCEPTED: canonical commit
accepted=1 rejected=2 committed=map[payment-commit-001:proposal-b]

[PROPOSAL] proposal-a
REJECTED: mutation already committed by proposal-b
accepted=1 rejected=3 committed=map[payment-commit-001:proposal-b]

=== CONSENSUS PROOF REPORT ===
session=session-xyz
epoch=1
authority=node-a
canonical_mutation=payment-commit-001
winner=proposal-b
accepted=1
rejected=3
invariant_violations=0
commit_once=true
fake_authority_rejected=true
stale_epoch_rejected=true
loser_rejected=true

VERDICT: CONSISTENT
Proof: exactly one proposal became canonical; all competing proposals failed closed