package history

////////////////////////////////////////////////////////////////

// Result is the result of a history walk.
type Result[E Event] struct {
	// Timeline is the timeline of events.
	Timeline []E

	// Subject is the subject that was used to process the events.
	Subject Subject[E]

	// Error is the error that occurred during the walk, if any.
	Error error
}

// NewResult creates a new Result with the given history, subject, and error.
//
// Parameters:
//   - history: The history of events.
//   - subject: The subject that was used to process the events.
//   - err: The error that occurred during the walk, if any.
//
// Returns:
//   - Result[E]: A new result with the given history, subject, and error.
func NewResult[E Event](history History[E], subject Subject[E], err error) Result[E] {
	timeline := TimelineOf(history)

	result := Result[E]{
		Timeline: timeline,
		Subject:  subject,
		Error:    err,
	}

	return result
}
