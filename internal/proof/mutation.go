package proof

type Event struct {
	Name       string
	SessionID  string
	MutationID string
	NodeID     string
	Authority  string
	Action     string
	Amount     int
}

func (e Event) CommitKey() string {
	return e.SessionID + ":" + e.MutationID
}