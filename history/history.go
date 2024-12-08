package history

import (
	"github.com/PlayerR9/go-evals/common"
)

// History is a history of events.
type History[E any] struct {
	// timeline is the sequence of events in the history.
	timeline []E

	// arrow is the current position in the history.
	arrow int
}

// AppendEvent creates a new history that appends the given event to the timeline.
// However, the current position in the history is not changed.
//
// Parameters:
//   - event: The event to append.
//
// Returns:
//   - History[E]: The new history.
func (h History[E]) AppendEvent(event E) History[E] {
	timeline := make([]E, len(h.timeline), len(h.timeline)+1)
	copy(timeline, h.timeline)

	timeline = append(timeline, event)

	new_history := History[E]{
		timeline: timeline,
		arrow:    h.arrow,
	}

	return new_history
}

// Walk walks the history and returns the next event. If the history is done, it returns
// an error. If the receiver is nil, it returns an error.
//
// Returns:
//   - E: The next event in the history.
//   - error: An error if the history could not be walked.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrExausted: If the history is done.
func (h *History[E]) Walk() (E, error) {
	if h == nil {
		return *new(E), common.ErrNilReceiver
	}

	if h.arrow == len(h.timeline) {
		return *new(E), ErrExausted
	}

	event := h.timeline[h.arrow]
	h.arrow++

	return event, nil
}

// Events returns a copy of the timeline.
//
// Returns:
//   - []E: A copy of the timeline.
func (h History[E]) Events() []E {
	if len(h.timeline) == 0 {
		return nil
	}

	slice := make([]E, len(h.timeline))
	copy(slice, h.timeline)

	return slice
}
