package history

import "errors"

var (
	// ErrExausted is the error returned when the history's iteration is done.
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"history is exhausted"
	ErrExausted error

	// ErrNoEvents occurs when there are no events to peek. This error can be
	// checked with the == operator.
	//
	// Format:
	// 	"no events to peek"
	ErrNoEvents error

	// ErrBreak is an error used to break out of a loop. This error can be checked
	// with the == operator.
	//
	// Format:
	// 	"break"
	ErrBreak error
)

func init() {
	ErrExausted = errors.New("history is exhausted")
	ErrNoEvents = errors.New("no events to peek")
	ErrBreak = errors.New("break")
}
