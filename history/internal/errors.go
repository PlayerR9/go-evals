package internal

import "errors"

var (
	// ErrEmptyQueue occurs when the queue is empty. This error can be checked with
	// the == operator.
	//
	// Format:
	// 	"queue is empty"
	ErrEmptyQueue error
)

func init() {
	ErrEmptyQueue = errors.New("queue is empty")
}
