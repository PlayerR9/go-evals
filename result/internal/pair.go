package internal

// Pair is a struct that represents a pair of results.
type Pair[T interface{ HasError() bool }] struct {
	// Results is a slice of Results.
	Results []T

	// IsValid is a boolean that indicates whether the pair is valid or not.
	IsValid bool
}

// NewPair creates a new Pair from a slice of results. If the slice contains at least one valid result,
// the returned Pair will be valid and contain the valid results. Otherwise, the returned Pair will be
// invalid and contain the invalid results.
//
// Parameters:
//   - inputs: The slice of results to create the Pair from.
//
// Returns:
//   - Pair[T]: The new Pair.
func NewPair[T interface{ HasError() bool }](inputs []T) Pair[T] {
	valids, invalids := Split(inputs)
	if len(valids) > 0 {
		return Pair[T]{
			Results: valids,
			IsValid: true,
		}
	} else {
		return Pair[T]{
			Results: invalids,
			IsValid: false,
		}
	}
}

// NewInvalidPair creates a new Pair from a slice of invalid results. The returned Pair will
// be invalid and contain the invalid results.
//
// Parameters:
//   - invalids: The slice of invalid results to create the Pair from.
//
// Returns:
//   - Pair[T]: The new Pair.
func NewInvalidPair[T interface{ HasError() bool }](invalids []T) Pair[T] {
	return Pair[T]{
		Results: invalids,
		IsValid: false,
	}
}

// RunFn is a function that takes a result and returns a pair of results.
//
// Parameters:
//   - elem: The result to evaluate.
//
// Returns:
//   - *Pair[T]: The pair of results.
//   - error: An error if the evaluation fails.
type RunFn[T interface{ HasError() bool }] func(elem T) (*Pair[T], error)
