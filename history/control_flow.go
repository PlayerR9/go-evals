package history

import "github.com/PlayerR9/go-evals/common"

// BodyFn is a function representing the body of a loop.
//
// Returns:
//   - error: An error if the loop fails.
//
// Errors:
//   - ErrBreak: If the loop should break.
//   - any error.
type BodyFn func() error

// Loop runs the given body function in a loop until it returns an error or
// ErrBreak.
//
// Parameters:
//   - bodyFn: The body function to run in the loop.
//
// Returns:
//   - error: An error if the loop failed.
//
// Errors:
//   - ErrNilParam: If the parameter is nil.
//   - any error: Implementation-specific errors.
func Loop(bodyFn BodyFn) error {
	if bodyFn == nil {
		err := common.NewErrNilParam("bodyFn")
		return err
	}

	for {
		err := bodyFn()
		if err == ErrBreak {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
