package history

import (
	"errors"

	"github.com/PlayerR9/go-evals/common"
)

////////////////////////////////////////////////////////////////

// History is a history of events.
type History[E Event] struct {
	// timeline is the list of events.
	timeline []E

	// arrow is the index of the current event.
	arrow uint
}

// TimelineOf returns a copy of the timeline of events.
//
// Parameters:
//   - h: The history to get the timeline of.
//
// Returns:
//   - []E: A copy of the timeline of events.
func TimelineOf[E Event](h History[E]) []E {
	if len(h.timeline) == 0 {
		return nil
	}

	timeline := make([]E, len(h.timeline))
	copy(timeline, h.timeline)

	return timeline
}

// ArrowOf returns the current index of the arrow in the timeline.
//
// Parameters:
//   - h: The history to get the arrow of.
//
// Returns:
//   - uint: The current index of the arrow in the timeline.
func ArrowOf[E Event](h History[E]) uint {
	return h.arrow
}

// AppendEvent creates a new History[E] that is a copy of the parameter but with the given
// event appended to the timeline.
//
// The new History[E] will have its arrow in the same position as the parameter.
//
// Parameters:
//   - h: The history to append the event to.
//   - event: The event to append to the timeline.
//
// Returns:
//   - History[E]: A new History[E] with the given event appended to the timeline.
func AppendEvent[E Event](h History[E], event E) History[E] {
	var timeline []E

	if len(h.timeline) > 0 {
		timeline = make([]E, len(h.timeline))
		copy(timeline, h.timeline)
	}

	timeline = append(timeline, event)

	result := History[E]{
		timeline: timeline,
		arrow:    h.arrow,
	}

	return result
}

// CurrentEventOf returns the event at the current position of the arrow, or a
// default-constructed event and false if the arrow is out of bounds.
//
// Parameters:
//   - h: The history to get the current event of.
//
// Returns:
//   - E: The event at the current position of the arrow, or a default-constructed
//     event if the arrow is out of bounds.
//   - bool: True if the arrow is in bounds, false otherwise.
func CurrentEventOf[E Event](h History[E]) (E, bool) {
	if h.arrow >= uint(len(h.timeline)) {
		return *new(E), false
	}

	return h.timeline[h.arrow], true
}

// Restart resets the arrow to the beginning of the timeline.
//
// Parameters:
//   - h: The history to restart.
//
// Returns:
//   - History[E]: A new History[E] with the arrow reset to the beginning of the timeline.
func Restart[E Event](h History[E]) History[E] {
	var timeline []E

	if len(h.timeline) > 0 {
		timeline = make([]E, len(h.timeline))
		copy(timeline, h.timeline)
	}

	history := History[E]{
		timeline: timeline,
		arrow:    0,
	}

	return history
}

// WalkForward moves the arrow forward and returns the event at the new position.
//
// The arrow is moved one position forward, and the event at the new position is
// returned. If the arrow is already at the end of the timeline, the function
// returns a default-constructed event and an error.
//
// Parameters:
//   - h: The history to walk.
//
// Returns:
//   - E: The event at the new position of the arrow, or a default-constructed event
//     if the arrow is at the end of the timeline.
//   - error: An error if the walk fails.
//
// Errors:
//   - common.ErrBadParam: If the history is nil.
//   - ErrEOT: If the end of the timeline is reached.
func WalkForward[E Event](h *History[E]) (E, error) {
	if h == nil {
		err := errors.New("parameter (h) must not be nil")
		return *new(E), err
	}

	if h.arrow >= uint(len(h.timeline)) {
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
// Parameters:
//   - h: The history to walk.
//
// Returns:
//   - E: The event at the new position of the arrow, or a default-constructed event
//     if the arrow is at the beginning of the timeline.
//   - error: An error if the walk fails.
//
// Errors:
//   - common.ErrBadParam: If the history is nil.
//   - ErrEOT: If the end of the timeline is reached.
func WalkBackward[E Event](h *History[E]) (E, error) {
	if h == nil {
		err := common.NewErrNilParam("h")
		return *new(E), err
	} else if h.arrow == 0 {
		return *new(E), ErrEOT
	}

	h.arrow--

	return h.timeline[h.arrow], nil
}
