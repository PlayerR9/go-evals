package history

// Subject is a subject that can be used to process events.
type Subject[E Event] interface {
	// HasError checks whether the subject has encountered non-fatal errors.
	//
	// Returns:
	//   - bool: True if the subject has encountered non-fatal errors, false otherwise.
	HasError() bool

	// ApplyEvent applies the given event to the subject.
	//
	// Parameters:
	//   - event: The event to apply.
	//
	// Returns:
	//   - error: An error if the subject cannot apply the event, nil otherwise.
	//
	// The error returned is for panic-level of errors. For any other error, make use
	// of the `HasError` method.
	ApplyEvent(event E) error

	// NextEvents returns the next events that can be applied to the subject.
	//
	// Returns:
	//   - []E: The next events that can be applied to the subject.
	//   - error: An error if the subject cannot return the next events, nil otherwise.
	//
	// The error returned is for panic-level of errors. For any other error, make use
	// of the `HasError` method.
	NextEvents() ([]E, error)
}

// alignOnce aligns the given subject with the given history by one step.
//
// The function walks the history forward and applies the next event to the subject.
// If the subject has errors after applying the event, the function returns an error.
//
// Parameters:
//   - subject: The subject to align.
//   - history: The history to align the subject with.
//
// Returns:
//   - History[E]: The updated history.
//   - error: An error if the subject has errors, nil otherwise.
//
// Errors:
//   - ErrEOT: If the end of the timeline is reached.
//   - ErrSubject: If the subject has errors.
//   - any error: The error returned by the subject's ApplyEvent method.
func alignOnce[E Event](subject Subject[E], history History[E]) (History[E], error) {
	// assert.Cond(subject != nil, "subject != nil")

	event, err := history.WalkForward()
	if err != nil {
		return history, err
	}

	err = subject.ApplyEvent(event)
	if err != nil {
		return history, err
	}

	ok := subject.HasError()
	if ok {
		return history, ErrSubject
	}

	return history, nil
}

// align repeatedly aligns the given subject with the given history until the end of the timeline.
//
// The function calls alignOnce to move the arrow forward and apply events to the subject.
// It continues this process until the end of the timeline is reached or an error occurs.
//
// Parameters:
//   - subject: The subject to align.
//   - history: The history with which to align the subject.
//
// Returns:
//   - History[E]: The final state of the history after alignment.
//   - error: An error if the alignment process fails.
//
// Errors:
//   - ErrEOT: If the end of the timeline is reached.
//   - any error: The error returned by the alignOnce function.
func align[E Event](subject Subject[E], history History[E]) (History[E], error) {
	// assert.Cond(subject != nil, "subject != nil")

	for {
		history, err := alignOnce(subject, history)
		if err == ErrEOT {
			break
		} else if err != nil {
			return history, err
		}
	}

	return history, nil
}

// walkOnce walks the history forward by one step and applies the event to the subject.
//
// This function retrieves the next event from the history and applies it to the subject.
//
// Parameters:
//   - subject: The subject to which the event will be applied.
//   - history: The history from which the next event is retrieved.
//
// Returns:
//   - History[E]: The updated history.
//   - error: An error if the walk or event application fails.
//
// Errors:
//   - ErrEOT: If the end of the timeline is reached.
//   - any error: The error returned by the subject's ApplyEvent method.
func walkOnce[E Event](subject Subject[E], history History[E]) (History[E], error) {
	// assert.Cond(subject != nil, "subject != nil")

	event, err := history.WalkForward()
	if err != nil {
		return history, err
	}

	err = subject.ApplyEvent(event)
	return history, err
}

// nextEvents retrieves the next possible histories by applying possible next events to the subject.
//
// This function retrieves the next events that can be applied to the subject and appends each event
// to the current history, creating a new history for each event. The new histories are returned
// in reverse order.
//
// Parameters:
//   - subject: The subject from which the next events are retrieved.
//   - history: The current history to which the events will be appended.
//
// Returns:
//   - []History[E]: A slice of new histories created by appending each possible next event to the current history.
//   - error: An error if the subject cannot return the next events, nil otherwise.
//
// Errors:
//   - any error: The error returned by the subject's NextEvents method.
func nextEvents[E Event](subject Subject[E], history History[E]) ([]History[E], error) {
	// assert.Cond(subject != nil, "subject != nil")

	nexts, err := subject.NextEvents()
	if err != nil {
		return nil, err
	}

	ok := subject.HasError()
	if ok || len(nexts) == 0 {
		return nil, nil
	}

	next_paths := make([]History[E], 0, len(nexts))

	for _, next := range nexts {
		path := history.AppendEvent(next)
		next_paths = append(next_paths, path)
	}

	return next_paths, nil
}
