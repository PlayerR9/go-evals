package matcher

// Automaton is an interface for an automaton.
type Automaton interface {
	// Close attempts to close the automaton.
	//
	// Returns:
	//   - error: An error if the close fails.
	//
	// Errors:
	//   - common.ErrNilReceiver: When the receiver is nil.
	//   - any other error: If the close fails, for whatever reason.
	Close() error

	// Reset resets the Automaton's internal state for reuse. Does
	// nothing if the receiver is nil.
	Reset()
}
