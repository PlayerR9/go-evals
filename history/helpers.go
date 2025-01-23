package history

import "reflect"

// TypeOf returns the type of the value as a string.
//
// Parameters:
//   - v: The value to get the type of.
//
// Returns:
//   - string: The type of the value as a string.
func TypeOf(v any) string {
	if v == nil {
		return "nil"
	}

	to := reflect.TypeOf(v)

	str := to.String()
	return str
}
