package result

import (
	"errors"
	"fmt"

	"github.com/PlayerR9/go-evals/common"
	"github.com/PlayerR9/go-evals/result/internal"
)

// Result is an interface for types that are used for handling multiple, chained
// results.
type Result interface {
	// HasError checks if the result has an error.
	//
	// Returns:
	// 	- bool: True if the result has an error, false otherwise.
	HasError() bool
}

// EvaluateFn is a function that evaluates an element and returns results.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - []T: The results of the evaluation.
//   - error: An error if the evaluation fails.
type EvaluateFn[T Result] func(elem T) ([]T, error)

// ApplyOnValidsFn is a function that processes elements and returns a slice of results.
//
// Parameters:
//   - elems: The elements to process.
//
// Returns:
//   - []T: The results of the processing.
//   - error: Nil if the processing is successful, an error if the processing fails.
type ApplyOnValidsFn[T Result] func(elems []T) ([]T, error)

// MakeRunFn is a function that returns a new RunFn that evaluates an element and applies a middle
// function to the valid results. If no valid results are found, it returns the invalid results with a
// boolean indicating success.
//
// Parameters:
//   - eval_fn: The EvaluateFn function that evaluates an element.
//   - apply_fn: The ApplyOnValidsFn function that processes elements.
//
// Returns:
//   - RunFn[T]: The RunFn function.
//   - error: An error if the evaluation function is nil.
func MakeRunFn[T Result](evalFn EvaluateFn[T], applyFn ApplyOnValidsFn[T]) (internal.RunFn[T], error) {
	if evalFn == nil {
		return nil, common.NewErrNilParam("evalFn")
	}

	var runFn internal.RunFn[T]

	if applyFn == nil {
		runFn = func(elem T) (*internal.Pair[T], error) {
			results, err := evalFn(elem)
			if err != nil {
				return nil, err
			} else if len(results) == 0 {
				return nil, nil
			}

			p := internal.NewPair(results)
			return &p, nil
		}
	} else {
		runFn = func(elem T) (*internal.Pair[T], error) {
			results, err := evalFn(elem)
			if err != nil {
				return nil, err
			} else if len(results) == 0 {
				return nil, nil
			}

			valids, invalids := internal.Split(results)
			if len(valids) == 0 {
				p := internal.NewInvalidPair(invalids)
				return &p, nil
			}

			res, err := applyFn(valids)
			if err == nil {
				p := internal.NewPair(res)
				return &p, nil
			} else if err != ErrInvalidResult {
				return nil, err
			}

			p := internal.NewInvalidPair(res)
			return &p, nil
		}
	}

	return runFn, nil
}

// MakeApplyFn creates a function that applies a RunFn to a slice of elements, aggregating valid
// and invalid results.
//
// Parameters:
//   - runFn: The RunFn function that processes an element.
//
// Returns:
//   - ApplyOnValidsFn[T]: A function that processes a slice of elements, returning valid results
//     and an error if any occur.
//   - error: An error if the runFn is nil.
//
// The returned function iterates over the slice of elements, applying runFn to each element.
// It aggregates valid results into one slice and invalid results into another. If any errors
// occur during processing, they are returned as a combined error. If valid results are present,
// they are returned without error; otherwise, invalid results are returned with an ErrInvalidResult.
func MakeApplyFn[T Result](runFn internal.RunFn[T]) (ApplyOnValidsFn[T], error) {
	if runFn == nil {
		return nil, common.NewErrNilParam("runFn")
	}

	fn := func(elems []T) ([]T, error) {
		if len(elems) == 0 {
			return nil, nil
		}

		var valid_sols []T
		invalid_sols := make([]T, 0)

		var errs []error

		for i, elem := range elems {
			p, err := runFn(elem)
			if err != nil {
				err := fmt.Errorf("index %d: %w", i, err)
				errs = append(errs, err)
			} else if p != nil {
				if p.IsValid {
					valid_sols = append(valid_sols, p.Results...)

					if invalid_sols != nil {
						clear(invalid_sols)
						invalid_sols = nil
					}
				} else if invalid_sols != nil {
					invalid_sols = append(invalid_sols, p.Results...)
				}
			}
		}

		err := errors.Join(errs...)
		if err != nil {
			return append(valid_sols, invalid_sols...), err
		} else if len(valid_sols) > 0 {
			return valid_sols, nil
		} else {
			return invalid_sols, ErrInvalidResult
		}
	}

	return fn, nil
}
