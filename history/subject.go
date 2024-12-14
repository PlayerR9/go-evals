package history

import (
	"errors"
	"fmt"
	"iter"

	assert "github.com/PlayerR9/go-verify"
)

// Subject is an interface representing an entity that can walk through a history.
type Subject[E Event] interface {
	// ApplyEvent applies an event to the subject.
	//
	// Parameters:
	//   - event: The event to apply.
	//
	// Returns:
	//   - error: An error if the event could not be applied.
	//
	// NOTES: Because the returned error causes the immediate stop of the history,
	// use it only for panic-level error handling. For any other error, use the
	// HasError method to signify something non-critical happened.
	ApplyEvent(event E) error

	// NextEvents returns the next events in the subject.
	//
	// Returns:
	//   - []E: The next events in the subject.
	//   - error: An error if the next events could not be returned.
	//
	// NOTES:
	// 	- Because the returned error causes the immediate stop of the history,
	// 	use it only for panic-level error handling. For any other error, use the
	// 	HasError method to signify something non-critical happened.
	NextEvents() ([]E, error)

	// HasError checks whether the subject has an error.
	//
	// Returns:
	//   - bool: True if the subject has an error, false otherwise.
	HasError() bool

	// GetError returns the error associated with the subject. However, this is mostly
	// used as a builder for the error and, as such, it always assume an error
	// has, indeed, occurred.
	//
	// Returns:
	//   - error: The error associated with the subject.
	GetError() error
}

// realign realigns the history with the subject by applying each event in the history
// to the subject using ApplyEvent method. It stops when an error occurs or when the
// history is fully walked. It asserts that the subject does not have an error after
// applying each event.
//
// Parameters:
//   - history: The history to realign.
//   - subject: The subject to apply the events to.
//
// Returns:
//   - error: An error if the history could not be realigned.
func realign[E Event](history *History[E], subject Subject[E]) error {
	assert.Cond(history != nil, "history != nil")
	assert.Cond(subject != nil, "subject != nil")

	for {
		event, err := history.Walk()
		if err != nil {
			break
		}

		err = subject.ApplyEvent(event)
		if err != nil {
			err := fmt.Errorf("while applying event: %w", err)
			return err
		}

		ok := subject.HasError()
		if ok {
			err := subject.GetError()
			err = fmt.Errorf("subject has an error: %w", err)
			return err
		}
	}

	return nil
}

// InitFn is a function that returns a subject and an error.
//
// Returns:
//   - Subject[E]: The subject.
//   - error: An error if the subject could not be created.
type InitFn[E Event] func() (Subject[E], error)

// MakeIter returns a sequence of all possible subjects that can be obtained by executing the given initialisation
// function and then applying events to the subject. The initialisation function is called at most once per
// unique subject. The sequence is ordered, with the first element being the result of executing the initialisation
// function once, and the rest being the result of applying events to the previous element. If the initialisation
// function returns an error, it is skipped.
//
// If the subject has an error at any point, it is skipped. If the subject has an error at the end, it is yielded
// at the end of the sequence.
//
// Parameters:
//   - init_fn: The initialisation function to execute.
//
// Returns:
//   - iter.Seq2[S, error]: A sequence of all possible subjects that can be obtained by executing the initialisation
//     function and then applying events to the subject. Never returns nil.
func MakeIter[E Event](init_fn InitFn[E]) iter.Seq2[Pair[E], error] {
	if init_fn == nil {
		init_fn = func() (Subject[E], error) {
			err := errors.New("no init function was provided")
			return nil, err
		}
	}

	fn := func(yield func(Pair[E], error) bool) {
		var history History[E]

		var e Evaluator[E]

		err := e.Enqueue(&history)
		assert.Err(err, "e.Enqueue()")

		var invalids []Pair[E]

		for {
			first, err := e.Dequeue()
			if err != nil {
				break
			}

			subject, err := init_fn()
			if err == nil && subject == nil {
				err = errors.New("subject is nil")
			}

			if err != nil {
				timeline := first.Timeline()
				pair := NewPair(subject, timeline)

				_ = yield(pair, err)
				return
			}

			history, err := e.executeUntil(subject)

			timeline := history.Timeline()
			pair := NewPair(subject, timeline)

			if err != nil {
				_ = yield(pair, err)
				return
			}

			if ok := subject.HasError(); ok {
				invalids = append(invalids, pair)
			} else {
				if ok := yield(pair, nil); !ok {
					return
				}
			}
		}

		for _, invalid := range invalids {
			if ok := yield(invalid, nil); !ok {
				return
			}
		}
	}

	return fn
}
