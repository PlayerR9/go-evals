package history

import (
	"github.com/PlayerR9/go-evals/common"
)

// History is a history of events.
type History[E Event] struct {
	// timeline is the list of events.
	timeline []E

	// arrow is the index of the current event.
	arrow uint
}

// Timeline returns a copy of the timeline of events.
//
// Returns:
//   - []E: A copy of the timeline of events.
func (h History[E]) Timeline() []E {
	if len(h.timeline) == 0 {
		return nil
	}

	timeline := make([]E, len(h.timeline))
	copy(timeline, h.timeline)

	return timeline
}

// Arrow returns the current index of the arrow in the timeline.
//
// Returns:
//   - uint: The current index of the arrow in the timeline.
func (h History[E]) Arrow() uint {
	return h.arrow
}

// CurrentEvent returns the event at the current position of the arrow, or a
// default-constructed event and false if the arrow is out of bounds.
//
// Returns:
//   - E: The event at the current position of the arrow, or a default-constructed
//     event if the arrow is out of bounds.
//   - bool: True if the arrow is in bounds, false otherwise.
func (h History[E]) CurrentEvent() (E, bool) {
	lenTimeline := uint(len(h.timeline))

	if h.arrow >= lenTimeline {
		return *new(E), false
	}

	return h.timeline[h.arrow], true
}

// Restart resets the arrow to the beginning of the timeline.
//
// Returns:
//   - error: An error if the history could not be reset.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (h *History[E]) Restart() error {
	if h == nil {
		return common.ErrNilReceiver
	}

	h.arrow = 0

	return nil
}

// WalkForward moves the arrow forward and returns the event at the new position.
//
// The arrow is moved one position forward, and the event at the new position is
// returned. If the arrow is already at the end of the timeline, the function
// returns a default-constructed event and an error.
//
// Returns:
//   - E: The event at the new position of the arrow, or a default-constructed event
//     if the arrow is at the end of the timeline.
//   - error: An error if the walk fails.
//
// Errors:
//   - common.ErrNilReceiver: If the history is nil.
//   - ErrEOT: If the end of the timeline is reached.
func (h *History[E]) WalkForward() (E, error) {
	if h == nil {
		return *new(E), common.ErrNilReceiver
	}

	lenTimeline := uint(len(h.timeline))
	if h.arrow >= lenTimeline {
		return *new(E), ErrEOT
	}

	event := h.timeline[h.arrow]
	h.arrow++

	return event, nil
}

// WalkBackward moves the arrow backward and returns the event at the new position.
//
// The arrow is moved one position backward, and the event at the new position is
// returned. If the arrow is already at the beginning of the timeline, the
// function returns a default-constructed event and an error.
//
// Returns:
//   - E: The event at the new position of the arrow, or a default-constructed event
//     if the arrow is at the beginning of the timeline.
//   - error: An error if the walk fails.
//
// Errors:
//   - common.ErrNilReceiver: If the history is nil.
//   - ErrEOT: If the end of the timeline is reached.
func (h *History[E]) WalkBackward() (E, error) {
	if h == nil {
		return *new(E), common.ErrNilReceiver
	} else if h.arrow == 0 {
		return *new(E), ErrEOT
	}

	h.arrow--

	return h.timeline[h.arrow], nil
}

// AppendEvent creates a new History[E] that is a copy of the parameter but with the given
// event appended to the timeline.
//
// The new History[E] will have its arrow in the same position as the parameter.
//
// Parameters:
//   - event: The event to append to the timeline.
//
// Returns:
//   - History[E]: A new History[E] with the given event appended to the timeline.
func (h History[E]) AppendEvent(event E) History[E] {
	var timeline []E

	if len(h.timeline) > 0 {
		timeline = make([]E, len(h.timeline), len(h.timeline)+1)
		copy(timeline, h.timeline)
	}

	timeline = append(timeline, event)

	result := History[E]{
		timeline: timeline,
		arrow:    h.arrow,
	}

	return result
}
