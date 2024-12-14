package history

// Pair is a pair of a solution and a history.
type Pair[E Event] struct {
	// Subject is the solution.
	Subject Subject[E]

	// History is the history.
	History []E
}

// NewPair creates a new pair.
//
// Parameters:
//   - solution: The solution.
//   - history: The history.
//
// Returns:
//   - Pair[E]: The new pair.
func NewPair[E Event](solution Subject[E], history []E) Pair[E] {
	return Pair[E]{
		Subject: solution,
		History: history,
	}
}

// GetError returns the error associated with the subject, if any.
//
// Returns:
//   - error: The error associated with the subject, or nil if the subject
//     has no error.
func (p Pair[E]) GetError() error {
	ok := p.Subject.HasError()
	if !ok {
		return nil
	}

	err := p.Subject.GetError()
	return err
}
