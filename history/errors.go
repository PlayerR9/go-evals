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

// ErrInvalidType occurs when the type of a value is invalid.
type ErrInvalidType struct {
	// Want is the expected type of the value.
	Want any

	// Got is the actual type of the value.
	Got any
}

// Error implements error.
func (e ErrInvalidType) Error() string {
	want_str := TypeOf(e.Want)
	got_str := TypeOf(e.Got)

	return "want " + want_str + ", got " + got_str
}

// NewErrInvalidType creates a new ErrInvalidType.
//
// Parameters:
//   - got: The actual type of the value.
//   - want: The expected type of the value.
//
// Returns:
//   - error: The newly created ErrInvalidType error. Never returns nil.
func NewErrInvalidType(got any, want any) error {
	e := &ErrInvalidType{
		Want: want,
		Got:  got,
	}

	return e
}
