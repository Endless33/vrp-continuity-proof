package proof

import "fmt"

func PrintReport(engine *Engine) {
	fmt.Println("=== PROOF REPORT ===")
	fmt.Printf("committed_mutations=%d\n", len(engine.Committed))
	fmt.Printf("duplicates_rejected=%d\n", engine.DuplicatesRejected)
	fmt.Printf("non_authority_rejected=%d\n", engine.AuthorityRejected)
	fmt.Printf("balance=%d\n", engine.Balance)
	fmt.Printf("invariant_violations=%d\n", engine.InvariantViolations)
	fmt.Println()
	fmt.Println("VERDICT:", engine.Verdict())

	if engine.Verdict() == "CONSISTENT" {
		fmt.Println("Proof: no logical mutation committed more than once")
	}
}