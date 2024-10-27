package internal

// Split takes a slice of results and splits it into two slices: one with valid
// results (i.e. results that have no error) and one with invalid results (i.e.
// results that have an error).
//
// Parameters:
//   - results: The slice of results to split.
//
// Returns:
//   - []T: The valid results.
//   - []T: The invalid results.
func Split[T interface{ HasError() bool }](results []T) ([]T, []T) {
	if len(results) == 0 {
		return nil, nil
	}

	var count_valids int

	for _, result := range results {
		ok := result.HasError()
		if !ok {
			count_valids++
		}
	}

	if count_valids == 0 {
		return nil, results
	}

	valids := make([]T, 0, count_valids)
	invalids := make([]T, 0, len(results)-count_valids)

	for _, result := range results {
		ok := result.HasError()
		if ok {
			invalids = append(invalids, result)
		} else {
			valids = append(valids, result)
		}
	}

	return valids, invalids
}
