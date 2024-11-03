package matcher

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrMatchDone occurs when the Matcher is complete. Readers are expected to return
	// this error as is and not wrap it as callers are expected to check for this error
	// with the == operator.
	ErrMatchDone error
)

func init() {
	ErrMatchDone = errors.New("match done")
}

// ErrNotAsExpected occurs when a string is not as expected.
type ErrNotAsExpected struct {
	// Quote if true, the strings will be quoted before being printed.
	Quote bool

	// Kind is the kind of the string that is not as expected.
	Kind string

	// Expecteds are the strings that were expecteds.
	Expecteds []any

	// Got is the actual string.
	Got any
}

// Error implements the error interface.
func (e ErrNotAsExpected) Error() string {
	var kind string

	if e.Kind != "" {
		kind = e.Kind + " to be "
	}

	var got string

	if e.Got == "" {
		got = "nothing"
	} else if e.Quote {
		got = strconv.Quote(fmt.Sprint(e.Got))
	} else {
		got = fmt.Sprint(e.Got)
	}

	var builder strings.Builder

	builder.WriteString("expected ")
	builder.WriteString(kind)

	if len(e.Expecteds) > 0 {
		elems := make([]string, 0, len(e.Expecteds))

		for _, elem := range e.Expecteds {
			str := fmt.Sprint(elem)
			elems = append(elems, str)
		}

		elems = Sort(elems)

		if e.Quote {
			for i := range elems {
				elems[i] = strconv.Quote(elems[i])
			}
		}

		builder.WriteString(EitherOrString(elems))
	} else {
		builder.WriteString("something")
	}

	builder.WriteString(", got ")
	builder.WriteString(got)

	return builder.String()
}

// NewErrNotAsExpected is a convenience function that creates a new ErrNotAsExpected error with
// the specified kind, got value, and expected values.
//
// See common.NewErrNotAsExpected for more information.
func NewErrNotAsExpected(quote bool, kind string, got any, expecteds ...any) error {
	return &ErrNotAsExpected{
		Quote:     quote,
		Kind:      kind,
		Expecteds: expecteds,
		Got:       got,
	}
}
