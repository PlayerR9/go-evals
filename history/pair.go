package history

import (
	"fmt"
	"iter"
)

// Pair is a pair of a solution and a history.
type Pair[E any, S Subject[E]] struct {
	// Subject is the solution.
	Subject S

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
//   - Pair[E, S]: The new pair.
func NewPair[E any, S Subject[E]](solution S, history []E) Pair[E, S] {
	return Pair[E, S]{
		Subject: solution,
		History: history,
	}
}

// GetError returns the error associated with the subject, if any.
//
// Returns:
//   - error: The error associated with the subject, or nil if the subject
//     has no error.
func (p Pair[E, S]) GetError() error {
	if !p.Subject.HasError() {
		return nil
	}

	return p.Subject.GetError()
}

// Emulate returns a sequence of shadow subjects created by applying the events
// in the history of the given pair to the shadow of its subject. The sequence
// is generated by yielding each shadow subject after applying an event.
//
// Parameters:
//   - p: The pair containing the subject and its history.
//
// Returns:
//   - iter.Seq[S]: A sequence of shadow subjects. Each shadow is yielded after
//     applying an event from the history. Never returns nil.
//
// Panics if applying an event to a shadow fails as it should never happen.
func Emulate[E any, S interface {
	Shadow() (S, error)

	Subject[E]
}](p Pair[E, S]) iter.Seq[S] {
	return func(yield func(S) bool) {
		shadow, err := p.Subject.Shadow()
		if err != nil {
			panic(fmt.Errorf("failed to create shadow: %w", err))
		}

		for _, event := range p.History {
			err := shadow.ApplyEvent(event)
			if err != nil {
				panic(fmt.Errorf("failed to apply event: %w", err))
			}

			if !yield(shadow) {
				return
			}
		}
	}
}