package history

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	"github.com/PlayerR9/go-evals/history/internal"
	assert "github.com/PlayerR9/go-verify"
)

// Subject is an interface representing an entity that can walk through a history.
type Subject[E any] interface {
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
func realign[E any](history *History[E], subject Subject[E]) error {
	assert.Cond(history != nil, "history != nil")
	assert.Cond(subject != nil, "subject != nil")

	for {
		event, err := history.Walk()
		if err != nil {
			break
		}

		err = subject.ApplyEvent(event)
		if err != nil {
			return fmt.Errorf("while applying event: %w", err)
		}

		ok := subject.HasError()
		if ok {
			err := subject.GetError()
			return fmt.Errorf("subject has an error: %w", err)
		}
	}

	return nil
}

// pushPaths appends all alternative paths that can be generated by applying the given events
// to the given history to the given slice of paths.
//
// Parameters:
//   - nexts: The events to apply to the given history.
//   - history: The history to apply the events to.
//   - all_paths: The slice of paths to append to.
//
// Returns:
//   - History[E]: The new history.
//
// Notes:
//   - The given slice of paths is modified in place.
//   - The given history is modified in place.
//   - The given slice of paths is sorted in reverse order of when the paths were generated.
func pushPaths[E any](nexts []E, history History[E], all_paths *internal.Queue[History[E]]) History[E] {
	assert.Cond(all_paths != nil, "all_paths != nil")
	assert.Cond(len(nexts) > 0, "len(nexts) > 0")

	if len(nexts) == 1 {
		history = history.AppendEvent(nexts[0])
		return history
	}

	paths := make([]History[E], 0, len(nexts))

	for _, next := range nexts {
		path := history.AppendEvent(next)
		paths = append(paths, path)
	}

	slices.Reverse(paths)

	for _, path := range paths[:len(paths)-1] {
		err := all_paths.Enqueue(path)
		assert.Err(err, "all_paths.Enqueue(path)")
	}

	history = paths[len(paths)-1]

	return history
}

// executeUntil executes the given history until it is done or an event is found that
// was not expected by the given subject. It modifies the given slice of paths in place.
//
// Parameters:
//   - all_paths: The slice of paths to append to.
//   - subject: The subject to execute with the history.
//
// Returns:
//   - error: An error if the subject is nil, or if the subject got an error or is done
//     before the history could be aligned.
func executeUntil[E any](all_paths *internal.Queue[History[E]], subject Subject[E]) (History[E], error) {
	assert.Cond(all_paths != nil, "all_paths != nil")
	assert.Cond(subject != nil, "subject != nil")

	history, err := all_paths.Dequeue()
	assert.Err(err, "all_paths.Dequeue()")

	err = realign(&history, subject)
	if err != nil {
		return history, err
	}

	for {
		nexts, err := subject.NextEvents()
		if err != nil {
			return history, err
		}

		if len(nexts) == 0 {
			break
		}

		history = pushPaths(nexts, history, all_paths)

		if ok := subject.HasError(); ok {
			break
		}

		event, err := history.Walk()
		assert.Err(err, "history.Walk()")

		err = subject.ApplyEvent(event)
		if err != nil {
			return history, err
		}

		if ok := subject.HasError(); ok {
			break
		}
	}

	return history, nil
}

// InitFn is a function that returns a subject and an error.
//
// Returns:
//   - Subject[E]: The subject.
//   - error: An error if the subject could not be created.
type InitFn[E any] func() (Subject[E], error)

// Evaluate returns a sequence of all possible subjects that can be obtained by executing the given initialisation
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
func Evaluate[E any](init_fn InitFn[E]) iter.Seq2[Pair[E], error] {
	if init_fn == nil {
		init_fn = func() (Subject[E], error) {
			err := errors.New("no init function was provided")
			return nil, err
		}
	}

	fn := func(yield func(Pair[E], error) bool) {
		var all_paths internal.Queue[History[E]]

		var history History[E]

		err := all_paths.Enqueue(history)
		assert.Err(err, "all_paths.Enqueue(history)")

		var invalids []Pair[E]

		for {
			ok := all_paths.IsEmpty()
			if ok {
				break
			}

			subject, err := init_fn()
			if err != nil {
				timeline := history.Timeline()
				pair := NewPair(subject, timeline)

				_ = yield(pair, err)
				return
			}

			if subject == nil {
				timeline := history.Timeline()
				pair := NewPair(subject, timeline)

				err = errors.New("subject is nil")

				_ = yield(pair, err)
				return
			}

			history, err = executeUntil(&all_paths, subject)

			timeline := history.Timeline()
			pair := NewPair(subject, timeline)

			if err != nil {
				_ = yield(pair, err)
				return
			}

			if ok := subject.HasError(); ok {
				invalids = append(invalids, pair)

				continue
			}

			ok = yield(pair, nil)
			if !ok {
				return
			}
		}

		for _, invalid := range invalids {
			ok := yield(invalid, nil)
			if !ok {
				return
			}
		}
	}

	return fn
}
