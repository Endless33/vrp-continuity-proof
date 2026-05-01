package proof

type Engine struct {
	Committed           map[string]bool
	CommitCount         map[string]int
	Balance             int
	DuplicatesRejected  int
	AuthorityRejected   int
	InvariantViolations int
	Trace               []string
}

func NewEngine() *Engine {
	return &Engine{
		Committed:   make(map[string]bool),
		CommitCount: make(map[string]int),
		Trace:       make([]string, 0),
	}
}

func (e *Engine) Apply(event Event) {
	key := event.CommitKey()

	if event.NodeID != event.Authority {
		e.AuthorityRejected++
		e.Trace = append(e.Trace, "REJECTED: non-authoritative")
		return
	}

	if e.Committed[key] {
		e.DuplicatesRejected++
		e.Trace = append(e.Trace, "REJECTED: duplicate mutation")
		return
	}

	e.Committed[key] = true
	e.CommitCount[key]++
	e.Balance += event.Amount

	e.Trace = append(e.Trace, "ACCEPTED: canonical commit")
}

func (e *Engine) Verify() {
	for _, count := range e.CommitCount {
		if count > 1 {
			e.InvariantViolations++
		}
	}
}

func (e *Engine) Verdict() string {
	if e.InvariantViolations == 0 {
		return "CONSISTENT"
	}

	return "BROKEN"
}