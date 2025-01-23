package history

import "errors"

var (
	// ErrEOT occurs when the end of the timeline is reached.
	//
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"end of timeline"
	ErrEOT error = errors.New("end of timeline")

	// ErrBreak occurs when the loop should break.
	//
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"break"
	ErrBreak error = errors.New("break")

	// ErrSubject occurs when the subject has an error.
	//
	// This error can be checked with the == operator.
	//
	// Format:
	// 	"subject has an error"
	ErrSubject error = errors.New("subject has an error")
)
