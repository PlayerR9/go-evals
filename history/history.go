package history

import (
	"fmt"
)

// History is a history of events.
type History[E any] struct {
	// timeline is the sequence of events in the history.
	timeline []E

	// arrow is the current position in the history.
	arrow int
}

// Validate implements the assert.Validater interface.
func (h History[E]) Validate() error {
	if h.arrow < 0 || h.arrow > len(h.timeline) {
		return fmt.Errorf("arrow is not in [%d, %d]", 0, len(h.timeline))
	}

	return nil
}

// AppendEvent creates a new history that appends the given event to the timeline.
// However, the current position in the history is not changed.
//
// Parameters:
//   - event: The event to append.
//
// Returns:
//   - *History[E]: The new history. Never returns nil.
func (h History[E]) AppendEvent(event E) *History[E] {
	timeline := make([]E, len(h.timeline), len(h.timeline)+1)
	copy(timeline, h.timeline)

	timeline = append(timeline, event)

	return &History[E]{
		timeline: timeline,
		arrow:    h.arrow,
	}
}

// Reset resets the history for reuse.
func (h *History[E]) Reset() {
	if h == nil {
		return
	}

	h.arrow = 0
}

// Walk returns the next event in the history and true if such an event exists,
// or the zero value of type E and false otherwise. The current position in the
// history is incremented if an event is returned.
//
// Returns:
//   - E: The next event in the history. Never returns nil.
//   - bool: True if an event exists, false otherwise.
func (h *History[E]) Walk() (E, bool) {
	if h == nil {
		return *new(E), false
	}

	// common.Validate(h)

	if h.arrow == len(h.timeline) {
		return *new(E), false
	}

	event := h.timeline[h.arrow]
	h.arrow++

	return event, true
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
