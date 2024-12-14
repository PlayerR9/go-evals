package history

import (
	"github.com/PlayerR9/go-evals/common"
)

// EventApplier is an interface for a subject that can walk through a history.
type EventApplier[E any] interface {
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
}

// ShadowOfFn is a function that returns a shadow of a subject and an error.
//
// Returns:
//   - EventApplier[E]: The shadow of the subject.
//   - error: An error if the shadow could not be created.
//
// In most cases, the shadow is just the subject itself in its initial state.
type ShadowOfFn[E any] func() (EventApplier[E], error)

// Simulate returns a sequence of shadow subjects created by applying the events
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
func Simulate[E any](shadowOfFn ShadowOfFn[E], timeline []E, visitFn func(shadow EventApplier[E]) error) error {
	if shadowOfFn == nil {
		return common.NewErrNilParam("shadowOfFn")
	} else if visitFn == nil {
		return common.NewErrNilParam("visitFn")
	}

	shadow, err := shadowOfFn(

	shadow, err := subject.Shadow()
	if err != nil {
		return err
	}

	for _, event := range timeline {
		err := shadow.ApplyEvent(event)
		if err != nil {
			return err
		}

		err = visitFn(shadow)
		if err == ErrBreak {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
