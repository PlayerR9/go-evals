package history

import (
	"fmt"
	"iter"
	"slices"
	// "github.com/PlayerR9/mygo-lib/common"
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
	for {
		event, err := history.Walk()
		if err != nil {
			break
		}

		err = subject.ApplyEvent(event)
		if err != nil {
			return err
		}

		ok := subject.HasError()
		if ok {
			return fmt.Errorf("subject has an error: %w", subject.GetError())
		}
	}

	return nil
}

// pushPaths appends all alternative paths that can be generated by applying the given events
// to the given history to the given slice of paths. It replaces the given history with the
// first path that was generated.
//
// Parameters:
//   - nexts: The events to apply to the given history.
//   - history: The history to apply the events to.
//   - all_paths: The slice of paths to append to.
//
// Notes:
//   - The given slice of paths is modified in place.
//   - The given history is modified in place.
//   - The given slice of paths is sorted in reverse order of when the paths were generated.
func pushPaths[E any](nexts []E, history *History[E], all_paths *[]History[E]) {
	switch len(nexts) {
	case 0:
		// Do nothing.
	case 1:
		*history = (*history).AppendEvent(nexts[0])
	default:
		paths := make([]History[E], 0, len(nexts))

		for _, next := range nexts {
			path := (*history).AppendEvent(next)
			paths = append(paths, path)
		}

		slices.Reverse(paths)

		*all_paths = append(*all_paths, paths[:len(paths)-1]...)
		*history = paths[len(paths)-1]
	}
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
func executeUntil[E any, S Subject[E]](all_paths *[]History[E], subject S) (History[E], error) {
	history := (*all_paths)[0]
	*all_paths = (*all_paths)[1:]

	err := realign(&history, subject)
	if err != nil {
		return history, err
	}

	event, err := history.Walk()
	if err == nil {
		err := subject.ApplyEvent(event)
		if err != nil {
			return history, err
		}
	}

	var nexts []E

	for err == nil && !subject.HasError() {
		nexts, err = subject.NextEvents()

		pushPaths(nexts, &history, all_paths)

		if err != nil {
			*all_paths = append(*all_paths, history)
			break
		} else if len(nexts) == 0 || subject.HasError() {
			break
		}

		event, _ := history.Walk()

		err = subject.ApplyEvent(event)
	}

	return history, err
}

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
func Evaluate[E any, S Subject[E]](init_fn func() (S, error)) iter.Seq2[Pair[E, S], error] {
	if init_fn == nil {
		init_fn = func() (S, error) {
			return *new(S), nil
		}
	}

	return func(yield func(Pair[E, S], error) bool) {
		var history History[E]

		all_paths := []History[E]{history}

		var invalids []Pair[E, S]

		for len(all_paths) > 0 {
			subject, err := init_fn()
			if err != nil {
				_ = yield(Pair[E, S]{
					Subject: subject,
					History: history.Events(),
				}, err)

				return
			}

			history, err = executeUntil(&all_paths, subject)
			pair := NewPair(subject, history.Events())

			if err != nil {
				_ = yield(pair, err)

				return
			}

			if subject.HasError() {
				invalids = append(invalids, pair)
			} else if !yield(pair, nil) {
				return
			}
		}

		for _, invalid := range invalids {
			if !yield(invalid, nil) {
				return
			}
		}
	}
}
