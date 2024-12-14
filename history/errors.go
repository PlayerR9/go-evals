package history

import "errors"

var (
	// ErrTimelineEnded is the error returned when the timeline's end is reached.
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"timeline has ended"
	ErrTimelineEnded error

	// ErrBreak is an error used to break out of a loop. This error can be checked
	// with the == operator.
	//
	// Format:
	// 	"break"
	ErrBreak error

	// ErrEmptyQueue occurs when the queue is empty. This error can be checked with
	// the == operator.
	//
	// Format:
	// 	"queue is empty"
	ErrEmptyQueue error
)

func init() {
	ErrTimelineEnded = errors.New("timeline has ended")
	ErrBreak = errors.New("break")
	ErrEmptyQueue = errors.New("queue is empty")
}
