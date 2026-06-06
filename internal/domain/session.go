package domain

// SessionState represents where a sync session is in its lifecycle.
type SessionState int

const (
	Idle SessionState = iota
	Connecting
	Negotiating // exchanging version vectors
	Syncing     // sending/receiving deltas
	Converged
)

// SessionEvent triggers a state transition.
type SessionEvent int

const (
	Connect SessionEvent = iota
	Connected
	Negotiated
	DeltasExchanged
	Confirmed
	Failed
	Disconnected
)

// Transition advances the sync session state machine. Pure function.
func Transition(current SessionState, event SessionEvent) SessionState {
	return Idle // TODO: state machine logic
}
