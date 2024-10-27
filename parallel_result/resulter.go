package result

import (
	"context"
	"sync"

	"github.com/PlayerR9/go-evals/common"
	"github.com/PlayerR9/go-evals/parallel_result/internal"
)

// Resulter is an interface for types that are used for handling multiple, chained
// results.
type Resulter interface {
	// HasError checks if the result has an error.
	//
	// Returns:
	// 	- bool: True if the result has an error, false otherwise.
	HasError() bool
}

// EvalResultFn is a function that evaluates an element and returns results. Errors returned
// by this function are treated as panic-level type of error that abrupts the evaluation, regardless
// of at which stage it was called. For non-abrupt errors, the resulter interface provides the
// HasError() method to catch them.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - []T: The results of the evaluation.
//   - error: An error if the evaluation fails.
type EvalResultFn[T Resulter] func(elem T) ([]T, error)

// ApplyFn is a function that processes elements and returns a slice of results along with a boolean
// indicating success. If no elements are specified, it returns an empty slice and false.
//
// Parameters:
//   - ctx: The context to use for the evaluation.
//   - elem: The elements to process.
//
// Returns:
//   - []T: The results of the processing.
//   - error: Nil if the processing is successful, an error if the processing fails.
//
// Errors:
//   - ErrInvalidResult: If the evaluation is successful but the result is invalid.
//   - any other error: When the evaluation fails due to internal failure.
type ApplyCtxFn[T Resulter] func(ctx context.Context, elem []T) ([]T, error)

// ParallelEval is a function that returns a new ApplyCtxFn that evaluates elements in parallel
// and applies a middle function to the valid results. If no valid results are found, it returns the
// invalid results with a boolean indicating success.
//
// Parameters:
//   - eval_fn: The EvalResultFn function that evaluates an element and returns results.
//   - mid_fn: The ApplyCtxFn function that processes elements and returns a slice of results along with a boolean
//     indicating success.
//
// Returns:
//   - ApplyCtxFn[T]: The ApplyCtxFn function that processes elements and returns a slice of results along with a boolean
//     indicating success.
//   - error: An error if the evaluation function or the middle function is nil.
func MakeBatchFn[T Resulter](eval_fn EvalResultFn[T], mid_fn ApplyCtxFn[T]) (RunErrFn[T], error) {
	if eval_fn == nil {
		return nil, common.NewErrNilParam("eval_fn")
	}

	batch_fn := func(ctx context.Context, ch chan internal.Pair[T], elem T) error {
		if ch == nil {
			return common.NewErrNilParam("ch")
		}

		results, err := eval_fn(elem)
		if err != nil {
			return err
		} else if len(results) == 0 {
			return nil
		}

		if mid_fn == nil {
			ch <- internal.NewPair(results)
			return nil
		}

		valids, invalids := internal.Split(results)
		if len(valids) == 0 {
			ch <- internal.NewInvalidPair(invalids)
			return nil
		}

		res, err := mid_fn(ctx, valids)
		if err == nil {
			ch <- internal.NewPair(res)
			return nil
		} else if err != ErrInvalidResult {
			return err
		}

		ch <- internal.NewInvalidPair(res)
		return nil
	}

	return batch_fn, nil
}

// resultListener listens to a channel of Pair[T] and separates the results into
// solutions and non-solutions.
//
// Parameters:
//   - ch: A channel of Pair[T] to listen to.
//   - sols: A pointer to a slice where valid results are appended.
//   - non_sols: A pointer to a slice where invalid results are appended, if provided.
//
// Returns:
//   - error: An error if either the sols or non_sols receiver is nil.
//
// The function reads from the channel and appends results from valid pairs to
// the sols slice. If a non_sols slice is provided, it appends results from
// invalid pairs to it. If a valid result is found, it clears the non_sols
// slice and sets it to nil.
func resultListener[T Resulter](ch <-chan internal.Pair[T], sols, non_sols *[]T) error {
	if ch == nil {
		return nil
	} else if sols == nil {
		return common.NewErrNilParam("sols")
	} else if non_sols == nil {
		return common.NewErrNilParam("non_sols")
	}

	for p := range ch {
		if p.IsValid {
			*sols = append(*sols, p.Results...)

			if non_sols != nil {
				clear(*non_sols)
				*non_sols = nil
				non_sols = nil
			}
		} else if non_sols != nil {
			*non_sols = append(*non_sols, p.Results...)
		}
	}

	return nil
}

// Evaluate is a function that returns a new ApplyCtxFn that evaluates elements in parallel
// and applies a middle function to the valid results. If no valid results are found, it returns the
// invalid results with a boolean indicating success.
//
// Parameters:
//   - eval_fn: The EvalResultFn function that evaluates an element and returns results.
//   - mid_fn: The ApplyCtxFn function that processes elements and returns a slice of results along with a boolean
//     indicating success.
//
// Returns:
//   - ApplyCtxFn[T]: The ApplyCtxFn function that processes elements and returns a slice of results along with a boolean
//     indicating success.
//   - error: An error if the evaluation function is nil.
func Evaluate[T Resulter](batch_fn RunErrFn[T]) ApplyCtxFn[T] {
	fn := func(parent context.Context, elems []T) ([]T, error) {
		if len(elems) == 0 {
			return nil, nil
		} else if batch_fn == nil {
			return nil, common.NewErrNilParam("batch_fn")
		}

		ch := make(chan internal.Pair[T], len(elems))

		var wg sync.WaitGroup
		wg.Add(1)

		var sols, non_sols []T

		go func() {
			defer wg.Done()

			_ = resultListener(ch, &sols, &non_sols)
		}()

		err := ExecuteBatch(parent, ch, elems, batch_fn)

		close(ch)

		wg.Wait()

		if err != nil {
			return append(sols, non_sols...), err
		} else if len(sols) > 0 {
			return sols, nil
		} else {
			return non_sols, ErrInvalidResult
		}
	}

	return fn
}
