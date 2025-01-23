package common

import "errors"

var (
	// ErrNilReceiver occurs when a method is called on a receiver that was not expected to be
	// nil.
	//
	// This error can be checked with the == operator.
	//
	// Format:
	//
	//		"receiver must not be nil"
	ErrNilReceiver error = errors.New("receiver must not be nil")
)

// ErrNilParam is an error that occurs when a parameter is nil.
type ErrNilParam struct {
	// ParamName is the name of the parameter that is nil.
	ParamName string
}

// Error returns a string representation of the error.
func (e ErrNilParam) Error() string {
	if e.ParamName == "" {
		return "parameter must not be nil"
	} else {
		return "parameter (" + e.ParamName + ") must not be nil"
	}
}

// NewErrNilParam creates a new ErrNilParam error.
//
// Parameters:
//   - param_name: The name of the parameter that is nil.
//
// Returns:
//   - error: The newly created ErrNilParam error. Never returns nil.
func NewErrNilParam(param_name string) error {
	e := &ErrNilParam{
		ParamName: param_name,
	}

	return e
}
