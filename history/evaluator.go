package history

import (
	"slices"

	"github.com/PlayerR9/go-evals/common"
	assert "github.com/PlayerR9/go-verify"
)

// Evaluator is a history evaluator.
type Evaluator[E Event] struct {
	// queue is the queue of histories to evaluate.
	queue []*History[E]
}

// Enqueue adds the given history to the queue of histories to evaluate.
//
// Parameters:
//   - history: The history to add to the queue.
//
// Returns:
//   - error: An error if the history could not be added to the queue.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - common.ErrBadParam: If the history is nil.
func (e *Evaluator[E]) Enqueue(history *History[E]) error {
	if e == nil {
		return common.ErrNilReceiver
	} else if history == nil {
		err := common.NewErrNilParam("history")
		return err
	}

	e.queue = append(e.queue, history)

	return nil
}

// Dequeue removes and returns the first history from the queue of histories to evaluate.
//
// Returns:
//   - *History[E]: The first history in the queue.
//   - error: An error if the queue is empty or the receiver is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrEmptyQueue: If the queue is empty.
func (e *Evaluator[E]) Dequeue() (*History[E], error) {
	if e == nil {
		return nil, common.ErrNilReceiver
	}

	ok := len(e.queue) == 0
	if ok {
		return nil, ErrEmptyQueue
	}

	history := e.queue[0]
	e.queue = e.queue[1:]

	return history, nil
}

// pushPaths appends all alternative paths that can be generated by applying the given events
// to the given history to the given slice of paths.
//
// Parameters:
//   - nexts: The events to apply to the given history.
//   - history: The history to apply the events to.
//
// Returns:
//   - History[E]: The new history.
//
// Notes:
//   - The given slice of paths is modified in place.
//   - The given history is modified in place.
//   - The given slice of paths is sorted in reverse order of when the paths were generated.
func (e *Evaluator[E]) pushPaths(nexts []E, history *History[E]) *History[E] {
	assert.Cond(e != nil, "e != nil")
	assert.Cond(len(nexts) > 0, "len(nexts) > 0")
	assert.Cond(history != nil, "history != nil")

	if len(nexts) == 1 {
		h := (*history)
		h = h.AppendEvent(nexts[0])
		return &h
	}

	paths := make([]*History[E], 0, len(nexts))

	for _, next := range nexts {
		path := history.AppendEvent(next)
		paths = append(paths, &path)
	}

	slices.Reverse(paths)

	for _, path := range paths[:len(paths)-1] {
		err := e.Enqueue(path)
		assert.Err(err, "e.Enqueue(&path)")
	}

	history = paths[len(paths)-1]

	return history
}

// executeUntil executes the given history until it is done or an event is found that
// was not expected by the given subject. It modifies the given slice of paths in place.
//
// Parameters:
//   - subject: The subject to execute with the history.
//
// Returns:
//   - error: An error if the subject is nil, or if the subject got an error or is done
//     before the history could be aligned.
func (e *Evaluator[E]) executeUntil(subject Subject[E]) (*History[E], error) {
	assert.Cond(e != nil, "e != nil")
	assert.Cond(subject != nil, "subject != nil")

	history, err := e.Dequeue()
	assert.Err(err, "e.Dequeue()")

	err = realign(history, subject)
	if err != nil {
		return nil, err
	}

	for {
		nexts, err := subject.NextEvents()
		if err != nil {
			return nil, err
		}

		if len(nexts) == 0 {
			break
		}

		history = e.pushPaths(nexts, history)

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
