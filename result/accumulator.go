package internal

import "github.com/PlayerR9/go-evals/common"

// Accumulator is an accumulator for results that can either be valid or invalid; giving
// priority to valid results.
//
// An empty accumulator can either be created with the `var acc Accumulator[E]` syntax or
// the `acc := new(Accumulator[E])` constructor.
type Accumulator[E any] struct {
	// results contains the results of the accumulator.
	results []E

	// is_valid indicates whether the accumulator contains valid results.
	is_valid bool
}

// Reset resets the accumulator back to its initial state; empty and invalid.
//
// Returns:
//   - error: An error if the accumulator cannot be reset.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (a *Accumulator[E]) Reset() error {
	if a == nil {
		return common.ErrNilReceiver
	}

	if len(a.results) > 0 {
		clear(a.results)
		a.results = nil
	}

	a.is_valid = false

	return nil
}

// AddValid adds a valid result to the accumulator. If the accumulator currently contains
// invalid results (i.e., `a.is_valid` is false), all of the results are cleared out and
// discarded. The result is then added to the accumulator and the accumulator is marked as
// valid.
//
// Parameters:
//   - elem: The valid result to add to the accumulator.
//
// Returns:
//   - error: An error if it cannot add the result to the accumulator.
//
// Errors:
//   - common.ErrNilReceiver: If the accumulator is nil.
func (a *Accumulator[E]) AddValid(elem E) error {
	if a == nil {
		return common.ErrNilReceiver
	}

	if !a.is_valid && len(a.results) > 0 {
		clear(a.results)
		a.results = nil
	}

	a.results = append(a.results, elem)
	a.is_valid = true

	return nil
}

// AddInvalid adds an invalid result to the accumulator. If the accumulator currently contains
// valid results (i.e., `a.is_valid` is true), the new result is discarded and this function
// returns nil.
//
// Parameters:
//   - elem: The invalid result to add to the accumulator.
//
// Returns:
//   - error: An error if it cannot add the result to the accumulator.
//
// Errors:
//   - common.ErrNilReceiver: If the accumulator is nil.
func (a *Accumulator[E]) AddInvalid(elem E) error {
	if a == nil {
		return common.ErrNilReceiver
	} else if a.is_valid {
		return nil
	}

	a.results = append(a.results, elem)

	return nil
}

// Results returns a slice of all the results in the accumulator.
//
// Returns:
//   - []E: A slice of all the results in the accumulator.
func (a Accumulator[E]) Results() []E {
	if len(a.results) == 0 {
		return nil
	}

	results := make([]E, len(a.results))
	copy(results, a.results)

	return results
}

// IsValid checks if all the results in the accumulator are valid.
//
// Returns:
//   - bool: True if all results in the accumulator are valid, false otherwise.
func (a Accumulator[E]) IsValid() bool {
	return a.is_valid
}
