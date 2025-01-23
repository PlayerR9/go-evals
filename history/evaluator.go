package history

import (
	"fmt"
	"iter"
	"slices"

	"github.com/PlayerR9/go-evals/common"
	evres "github.com/PlayerR9/go-evals/result"
	assert "github.com/PlayerR9/go-verify"
)

////////////////////////////////////////////////////////////////

// InitFn is a type of function that is used to initialize a subject.
//
// Returns:
//   - Subject[E]: The initial state of the subject.
//   - error: An error if the initialization fails.
type InitFn[E Event] func() (Subject[E], error)

// Evaluator is an evaluator for event histories.
type Evaluator[E Event] struct {
	// initFn is a function that returns the initial state of the subject.
	initFn InitFn[E]

	// paths is a stack of event histories.
	paths []History[E]
}

// NewEvaluator creates a new Evaluator.
//
// Parameters:
//   - initFn: A function that returns the initial state of the subject.
//
// Returns:
//   - *Evaluator[E]: The newly created evaluator.
//   - error: An error if the evaluator could not be created.
//
// Errors:
//   - common.ErrNilParam: If the parameter is nil.
func NewEvaluator[S Subject[E], E Event](initFn func() (S, error)) (*Evaluator[E], error) {
	if initFn == nil {
		err := common.NewErrNilParam("initFn")
		return nil, err
	}

	fn := func() (Subject[E], error) {
		subject, err := initFn()
		if err != nil {
			return nil, err
		}

		return subject, nil
	}

	eval := &Evaluator[E]{
		initFn: fn,
		paths:  nil,
	}

	return eval, nil
}

// pop removes and returns the top history from the stack of event histories.
//
// The function removes the last history from the paths stack and returns it.
// If the stack is empty, the function returns an empty History and an error.
//
// Returns:
//   - History[E]: The top history from the stack.
//   - error: An error if the stack is empty.
func (e *Evaluator[E]) pop() (History[E], bool) {
	assert.NotNil(e, "e")

	if len(e.paths) == 0 {
		return History[E]{}, false
	}

	history := e.paths[len(e.paths)-1]
	e.paths = e.paths[:len(e.paths)-1]

	return history, true
}

// push pushes a slice of elements onto the stack in reverse order.
//
// Parameters:
//   - elems: The slice of elements to be pushed. If empty, no elements are pushed.
func (e *Evaluator[E]) push(elems []History[E]) {
	assert.NotNil(e, "e")

	if len(elems) == 0 {
		return
	}

	slices.Reverse(elems)

	e.paths = append(e.paths, elems...)
}

// applyOnce performs one iteration of the algorithm to find all possible histories from the current
// state of the subject and the current history.
//
// The function does the following:
//  1. Retrieves the next possible events from the subject.
//  2. Creates a new history for each of the next possible events by appending the event to the
//     current history.
//  3. Chooses the last new history and walks the history forward by one step.
//  4. Repeats steps 1-3 until there are no more next possible events or the subject has errors.
//
// Parameters:
//   - subject: The subject from which the next possible events are retrieved.
//   - history: The current history.
//   - paths: A pointer to a slice of histories that will be populated with the new histories
//     created by appending each possible next event to the current history.
//
// Returns:
//   - History[E]: The final history after all iterations.
//   - error: An error if any error occurs during the iteration, ErrContinue if the iteration should
//     continue, nil otherwise.
//
// Errors:
//   - ErrEOT: If the end of the timeline is reached.
//   - ErrBreak: If the subject has errors or there are no more next possible events.
//   - any error: The error returned by the subject's NextEvents method or the subject's ApplyEvent
//     method.
func (e *Evaluator[E]) applyOnce(subject Subject[E], history History[E]) (History[E], error) {
	assert.NotNil(e, "e")
	assert.Cond(subject != nil, "subject != nil")

	nexts, err := nextEvents(subject, history)
	if err != nil {
		return history, err
	} else if len(nexts) == 0 {
		return history, ErrBreak
	}

	history = nexts[len(nexts)-1]

	nexts = nexts[:len(nexts)-1]
	e.push(nexts)

	history, err = walkOnce(subject, history)
	if err != nil {
		return history, err
	} else if subject.HasError() {
		return history, ErrBreak
	}

	return history, nil
}

// apply returns a sequence of Results that contain all possible histories from the
// initial state of the subject created by initFn and the empty history.
//
// The sequence is lazily generated by walking the history forward by one step at a
// time and applying the next event to the subject. When the subject has errors, the
// sequence will yield the history and the subject at the point of error. When the
// subject has no more next events, the sequence will yield the final history and
// the subject.
//
// The sequence will not yield any results if the subject has errors after the
// initial alignment.
//
// Parameters:
//   - initFn: A function that creates a new subject.
//
// Returns:
//   - iter.Seq[Result[E]]: A sequence of Results that contain all possible histories
//     from the initial state of the subject created by initFn and the empty history.
func (e *Evaluator[E]) apply() iter.Seq[Result[E]] {
	assert.NotNil(e, "e")

	fn := func(yield func(Result[E]) bool) {
		var h History[E]

		e := Evaluator[E]{
			initFn: e.initFn,
		}

		e.push([]History[E]{h})

		var invalids []Result[E]

		for {
			top, ok := e.pop()
			if !ok {
				break
			}

			top = Restart(top)

			subject, err := e.initFn()
			if err != nil {
				_ = yield(NewResult(top, subject, fmt.Errorf("initFn() failed: %w", err)))
				return
			}

			top, err = align(subject, top)
			if err != nil {
				_ = yield(NewResult(top, subject, fmt.Errorf("align() failed: %w", err)))
				return
			}

			fn := func() error {
				res, err := e.applyOnce(subject, top)
				top = res
				return err
			}

			for {
				err := fn()
				if err == ErrBreak {
					break
				} else if err != nil {
					_ = yield(NewResult(top, subject, fmt.Errorf("doOne() failed: %w", err)))
					return
				}
			}

			r := NewResult(top, subject, nil)

			if subject.HasError() {
				invalids = append(invalids, r)
			} else if !yield(r) {
				return
			}
		}

		for _, r := range invalids {
			if !yield(r) {
				return
			}
		}
	}

	return fn
}

// AsSeq returns a sequence of Results that contain all possible histories from the
// initial state of the subject created by initFn and the empty history.
//
// The sequence is lazily generated by walking the history forward by one step at a
// time and applying the next event to the subject. When the subject has errors, the
// sequence will yield the history and the subject at the point of error. When the
// subject has no more next events, the sequence will yield the final history and
// the subject.
//
// The sequence will not yield any results if the subject has errors after the
// initial alignment.
//
// Parameters:
//   - initFn: A function that creates a new subject.
//
// Returns:
//   - iter.Seq[Result[E]]: A sequence of Results that contain all possible histories
//     from the initial state of the subject created by initFn and the empty history.
//
// When initFn is nil, it defaults to a function that returns a new subject and an
// error.
func (e *Evaluator[E]) AsSeq() iter.Seq[Result[E]] {
	if e == nil {
		return func(yield func(Result[E]) bool) {}
	}

	seq := e.apply()
	return seq
}

// execute executes the sequence generated by apply and returns all valid results.
//
// If any of the results contain an error, the function returns the first error
// encountered. If any of the results contain an invalid subject, the function
// returns all invalid results. If all results are valid, the function returns all
// valid results.
//
// The function does not return any results if the subject has errors after the
// initial alignment.
//
// Parameters:
//   - initFn: A function that creates a new subject.
//
// Returns:
//   - []Result[E]: A slice of Results that contain all valid histories from the
//     initial state of the subject created by initFn and the empty history.
//   - error: The first error encountered during the execution of the sequence, or
//     nil if no errors were encountered.
func (e *Evaluator[E]) execute() ([]Result[E], error) {
	assert.Cond(e != nil, "e != nil")

	seq := e.apply()

	var builder evres.Accumulator[Result[E]]
	defer builder.Reset()

	for res := range seq {
		if res.Error != nil {
			return builder.Results(), res.Error
		}

		ok := res.Subject.HasError()

		if !ok {
			err := builder.AddValid(res)
			assert.Err(err, "builder.AddValid(res)")
		} else {
			err := builder.AddInvalid(res)
			assert.Err(err, "builder.AddInvalid(res)")
		}
	}

	return builder.Results(), nil
}

// Execute executes a sequence of Results that contain all valid histories from the
// initial state of the subject created by initFn and the empty history.
//
// Returns:
//   - []Result[E]: A slice of Results that contain all valid histories from the
//     initial state of the subject created by initFn and the empty history.
//   - error: The first error encountered during the execution of the sequence, or
//     nil if no errors were encountered.
func (e *Evaluator[E]) Execute() ([]Result[E], error) {
	if e == nil {
		return nil, common.ErrNilReceiver
	}

	results, err := e.execute()
	return results, err
}
