package result

import "errors"

var (
	// ErrInvalidResult occurs when a result is invalid. This is used as a flag to signify that, the
	// resulting slice of results indicate the invalid results. If the error is different than this
	// one using == operator, then the error is an evaluation error rather than a result error.
	//
	// Format:
	//
	//   "invalid result"
	ErrInvalidResult error
)

func init() {
	ErrInvalidResult = errors.New("invalid result")
}
