package matcher

import "github.com/PlayerR9/go-evals/common"

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

	// Matched returns the matched characters as a slice of runes. The returned slice is
	// a copy and is valid until the next call to Reset or Match.
	//
	// Returns:
	//   - []I: The matched characters.
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
