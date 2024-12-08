package history

// Resetter is an interface for a resetter.
type Resetter interface {
	// Reset resets the Resetter's internal state for reuse.
	//
	// Returns:
	//   - error: An error if the reset fails.
	//
	// Errors:
	//   - common.ErrNilReceiver: When the receiver is nil.
	//   - any other error: If the reset fails, for whatever reason.
	Reset() error
}
