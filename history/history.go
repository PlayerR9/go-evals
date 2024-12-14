package history

import (
	"github.com/PlayerR9/go-evals/common"
)

// History is a history of events.
type History[E Event] struct {
	// timeline is the sequence of events in the history.
	timeline []E

	// arrow is the current position in the history.
	arrow uint
}

// Timeline returns a copy of the timeline of events in the history.
//
// Returns:
//   - []E: A copy of the timeline or nil if empty.
func (h History[E]) Timeline() []E {
	if len(h.timeline) == 0 {
		return nil
	}

	timeline := make([]E, len(h.timeline))
	copy(timeline, h.timeline)

	return timeline
}

// Arrow returns the current position in the history. The position is a 0-indexed
// offset into the timeline of events.
//
// Returns:
//   - uint: The current position in the history.
func (h History[E]) Arrow() uint {
	return h.arrow
}

// PeekEvent returns the current event in the history without advancing the position.
//
// Returns:
//   - E: The current event in the history.
//   - error: An error if there is no current event.
//
// Errors:
//   - ErrHistoryDone: If there is no current event.
func (h History[E]) PeekEvent() (E, error) {
	if h.arrow >= uint(len(h.timeline)) {
		return *new(E), ErrTimelineEnded
	}

	event := h.timeline[h.arrow]

	return event, nil
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

// Walk walks the history and returns the next event.
//
// Returns:
//   - E: The next event in the history.
//   - error: An error if the history could not be walked.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrHistoryDone: If the history is done.
func (h *History[E]) Walk() (E, error) {
	if h == nil {
		return *new(E), common.ErrNilReceiver
	}

	if h.arrow >= uint(len(h.timeline)) {
		return *new(E), ErrTimelineEnded
	}

	event := h.timeline[h.arrow]
	h.arrow++

	return event, nil
}
