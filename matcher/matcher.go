package matcher

import (
	"errors"

	"github.com/PlayerR9/go-evals/common"
	"github.com/PlayerR9/go-evals/rank"
)

// Matcher is an interface for matching elements.
type Matcher[I comparable] interface {
	// Match attempts to match the given element with the Matcher.
	//
	// Parameters:
	//   - elem: The element to be matched.
	//
	// Returns:
	//   - error: An error if the match fails or if it is complete.
	//
	// Errors:
	//   - ErrMatchDone: The Matcher is complete.
	//   - any other error: If the match fails, for whatever reason.
	Match(elem I) error

	// Matched returns the matched elements as a slice of elements. The returned slice is
	// a copy and is valid until the next call to Reset or Match.
	//
	// Returns:
	//   - []I: The matched elements.
	Matched() []I

	Automaton
}

// Execute runs the given Matcher on the given slice of elements, and returns the
// matched elements. If the Matcher completes early (i.e., returns ErrMatchDone),
// the execution will stop early. If the Matcher fails to match, an error will be
// returned. If the Matcher completes normally, the matched elements will be
// returned, and the Matcher will be reset.
//
// Parameters:
//   - m: The Matcher to be executed.
//   - slice: The slice of elements to be matched.
//
// Returns:
//   - []I: The matched elements.
//   - error: An error if the match fails, or if the Matcher completes early.
func Execute[I comparable](m Matcher[I], slice []I) ([]I, error) {
	if len(slice) == 0 {
		return nil, nil
	} else if m == nil {
		return nil, common.NewErrNilParam("m")
	}

	defer m.Reset()

	early_exit := false

	for i := 0; i < len(slice) && !early_exit; i++ {
		err := m.Match(slice[i])
		if err == nil {
			continue
		}

		if err == ErrMatchDone {
			early_exit = true
		} else {
			return m.Matched(), err
		}
	}

	if early_exit {
		return m.Matched(), nil
	}

	err := m.Close()
	return m.Matched(), err
}

// Pair is a pair of indices and matched elements.
type Pair[I comparable] struct {
	// Idx is the index of the Matcher that successfully matched the elements.
	Idx int

	// Matched is the matched elements.
	Matched []I
}

// Match uses the provided matchers to process a sequence of elements. It attempts
// to match each element with the matchers specified by the indices.
//
// Parameters:
//   - matchers: A slice of Matcher instances to be applied to the elements.
//   - indices: The initial indices of the matchers to be used for matching.
//   - elems: The elements to be matched.
//
// Returns:
//   - []Pair[T]: A slice of indices-matched pairs indicating the matchers that successfully
//     matched the elements.
//   - error: An error if any matcher fails or if there are issues completing the
//     matching process.
//
// The resulting slice of pairs is sorted from the one that matched the longest
// to the one that matched the shortest.
func Match[I comparable](matchers []Matcher[I], indices []int, elems []I) ([]Pair[I], error) {
	defer func() {
		for _, m := range matchers {
			m.Reset()
		}
	}()

	eos := rank.NewErrRorSol[int]()
	var level int

	for _, elem := range elems {
		if len(indices) == 0 {
			break
		}

		var top int

		for _, idx := range indices {
			m := matchers[idx]

			err := m.Match(elem)
			if err == nil {
				indices[top] = idx
				top++
			} else if err == ErrMatchDone {
				_ = eos.AddSol(level, idx)
			} else {
				_ = eos.AddErr(level, err)
			}
		}

		indices = indices[:top:top]
		level++
	}

	for _, idx := range indices {
		m := matchers[idx]

		err := m.Close()
		if err != nil {
			_ = eos.AddErr(level, err)
		} else {
			_ = eos.AddSol(level, idx)
		}
	}

	if eos.HasError() {
		return nil, errors.Join(eos.Errors()...)
	}

	sols := eos.Sols()

	results := make([]Pair[I], 0, len(sols))

	for _, sol := range sols {
		m := matchers[sol]

		p := Pair[I]{
			Idx:     sol,
			Matched: m.Matched(),
		}

		results = append(results, p)
	}

	return results, nil
}
