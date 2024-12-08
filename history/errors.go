package history

import "errors"

var (
	// ErrExausted is the error returned when the history's iteration is done.
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"history is exhausted"
	ErrExausted error
)

func init() {
	ErrExausted = errors.New("history is exhausted")
}
