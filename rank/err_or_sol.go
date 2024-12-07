package rank

import "github.com/PlayerR9/go-evals/common"

// ErrOrSol is a data structure that holds not only a list of solutions, but also a list of errors.
// If at least one solution is given, the list of errors is ignored and emptied.
//
// An empty ErrOrSol can either be created with the `var eos ErrOrSol[T]` syntax or with the
// `new(ErrOrSol[T])` constructor.
type ErrOrSol[T any] struct {
	// errs is the list of errors.
	errs []error

	// sols is the list of solutions.
	sols []T
}

// Size returns the size of the EOS.
//
// Returns:
//   - int: The size of the EOS. If there are solutions,
//     their size is returned, otherwise, the size of
//     the errors is returned. Never negative.
func (eos ErrOrSol[T]) Size() int {
	if len(eos.sols) > 0 {
		return len(eos.sols)
	} else {
		return len(eos.errs)
	}
}

// IsEmpty checks whether the ErrOrSol EOS is empty.
//
// Returns:
//   - bool: True if both the solutions and errors are empty, false otherwise.
func (eos ErrOrSol[T]) IsEmpty() bool {
	if len(eos.sols) > 0 {
		return false
	}

	return len(eos.errs) == 0
}

// Reset resets the EOS for reuse.
func (eos *ErrOrSol[T]) Reset() {
	if eos == nil {
		return
	}

	if eos.errs != nil {
		clear(eos.errs)
		eos.errs = make([]error, 0)
	}

	if eos.sols != nil {
		clear(eos.sols)
		eos.sols = nil
	}
}

// AddSol adds a solution to the EOS. If at least one solution is added,
// the errors are discarded and ignored.
//
// Parameters:
//   - elem: The solution to add.
//
// Returns:
//   - error: An error if the receiver is nil.
func (eos *ErrOrSol[T]) AddSol(elem T) error {
	if eos == nil {
		return common.ErrNilReceiver
	}

	if len(eos.errs) > 0 {
		clear(eos.errs)
		eos.errs = nil
	}

	eos.sols = append(eos.sols, elem)

	return nil
}

// AddErr adds an error to the EOS.
//
// Parameters:
//   - err: The error to add.
//
// Returns:
//   - error: An error if the receiver is nil.
//
// Behaviors:
//   - If the error is nil, it is ignored.
//   - If at least a solution has been added, the error is ignored.
func (eos *ErrOrSol[T]) AddErr(err error) error {
	if err == nil {
		return nil
	} else if eos == nil {
		return common.ErrNilReceiver
	}

	if len(eos.sols) > 0 {
		return nil
	}

	eos.errs = append(eos.errs, err)

	return nil
}

// HasError checks whether the EOS has an error.
//
// Returns:
//   - bool: True if the EOS has an error, false otherwise.
func (eos ErrOrSol[T]) HasError() bool {
	return len(eos.errs) > 0
}

// Errors returns the list of errors in descending order of rank.
//
// Returns:
//   - []error: The list of errors. Nil if there are no errors.
func (eos ErrOrSol[T]) Errors() []error {
	if len(eos.errs) == 0 {
		return nil
	}

	errs := make([]error, len(eos.errs))
	copy(errs, eos.errs)

	return errs
}

// Sols returns the list of solutions in descending order of rank.
//
// Returns:
//   - []T: The list of solutions. Nil if there are no solutions.
func (eos ErrOrSol[T]) Sols() []T {
	if len(eos.sols) == 0 {
		return nil
	}

	sols := make([]T, len(eos.sols))
	copy(sols, eos.sols)

	return sols
}
