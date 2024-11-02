package matcher

import (
	"slices"
	"strings"
)

// Predicate is a function that takes an element and returns true if it matches.
//
// Parameters:
//   - elem: The element to match.
//
// Returns:
//   - bool: True if the element matches, false otherwise.
type Predicate[I comparable] func(elem I) bool

// RejectNils returns the number of elements in the slice that are nil. The function
// modifies the slice by removing all nil elements and returns the number of elements
// that were removed. If the resulting slice is empty, the function sets the slice to nil.
//
// The function is useful when a nil slice or nil elements are not expected and should
// be removed.
//
// Parameters:
//   - slice: The slice to process.
//
// Returns:
//   - int: The number of elements removed. Never returns a negative number.
func RejectNils[I comparable](slice *[]Matcher[I]) int {
	if slice == nil || len(*slice) == 0 {
		return 0
	}

	n := len(*slice)

	var top int

	for _, s := range *slice {
		if s != nil {
			(*slice)[top] = s
			top++
		}
	}

	if top == 0 {
		clear(*slice)
		*slice = nil
	} else {
		clear((*slice)[top:])
		*slice = (*slice)[:top:top]
	}

	return n - top
}

// EitherOrString is a function that returns a string representation of a slice
// of strings. Empty strings are ignored.
//
// Parameters:
//   - values: The values to convert to a string.
//
// Returns:
//   - string: The string representation.
//
// Example:
//
//	EitherOrString([]string{"a", "b", "c"}) // "either a, b, or c"
func EitherOrString(elems []string) string {
	var str string

	switch len(elems) {
	case 0:
		// Do nothing
	case 1:
		str = elems[0]
	case 2:
		str = "either " + elems[0] + " or " + elems[1]
	default:
		str = "either " + strings.Join(elems[:len(elems)-1], ", ") + ", or " + elems[len(elems)-1]
	}

	return str
}

// Sort removes duplicate strings from the given slice and sorts it in-place.
//
// This function only works on a slice of strings.
//
// Parameters:
//   - elems: The slice to sort.
//
// Example:
//
//	elems := []string{"d", "b", "a", "c", "b", "c"}
//	common.Sort(&elems) // elems is now []string{"a", "b", "c", "d"}
func Sort(elems *[]string) {
	if elems == nil || len(*elems) == 0 {
		return
	}

	unique := make([]string, 0, len(*elems))

	for _, elem := range *elems {
		pos, ok := slices.BinarySearch(*elems, elem)
		if ok {
			continue
		}

		unique = slices.Insert(unique, pos, elem)
	}

	size := len(unique)

	unique = unique[:size:size]

	clear((*elems)[size:])
	*elems = unique
}
