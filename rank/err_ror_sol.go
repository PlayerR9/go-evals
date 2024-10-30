package rank

import "github.com/PlayerR9/go-evals/common"

// ErrRorSol is a data structure that holds not only a list of solutions
// according to their rank, but also a list of errors according to their
// their own separated ranking. If at least one solution is given, the
// list of errors is ignored and emptied.
type ErrRorSol[T any] struct {
	// errs is the list of errors.
	errs *Rank[error]

	// sols is the list of solutions.
	sols *Rank[T]
}

// NewErrRorSol creates and returns a new instance of ErrRorSol.
// The returned instance is initialized with an empty list of errors
// and no solutions.
//
// Returns:
//   - *ErrRorSol[T]: A new instance of ErrRorSol with a dedicated
//     ranking for errors, ready to use. Never returns nil.
func NewErrRorSol[T any]() *ErrRorSol[T] {
	return &ErrRorSol[T]{
		errs: new(Rank[error]),
		sols: nil,
	}
}

// Size returns the size of the EOS.
//
// Returns:
//   - int: The size of the EOS. If there are solutions,
//     their size is returned, otherwise, the size of
//     the errors is returned. Never negative.
func (eos ErrRorSol[T]) Size() int {
	if eos.errs == nil {
		return eos.sols.Size()
	} else {
		return eos.errs.Size()
	}
}

// IsEmpty checks whether the ErrRorSol EOS is empty.
//
// Returns:
//   - bool: True if both the solutions and errors are empty, false otherwise.
func (eos ErrRorSol[T]) IsEmpty() bool {
	if eos.errs == nil {
		return eos.sols.IsEmpty()
	} else {
		return eos.errs.IsEmpty()
	}
}

// Reset resets the EOS for reuse.
func (eos *ErrRorSol[T]) Reset() {
	if eos == nil {
		return
	}

	if eos.errs != nil {
		eos.errs.Reset()
		eos.errs = new(Rank[error])
	}

	if eos.sols != nil {
		eos.sols.Reset()
		eos.sols = nil
	}
}

// AddSol adds a solution to the EOS. If at least one solution is added,
// the errors are discarded and ignored.
//
// Parameters:
//   - rank: The level of the solution.
//   - elem: The solution to add.
//
// Returns:
//   - error: An error if the receiver is nil.
func (eos *ErrRorSol[T]) AddSol(rank int, elem T) error {
	if eos == nil {
		return common.ErrNilReceiver
	}

	if eos.errs != nil {
		eos.errs.Reset()
		eos.errs = nil
	}

	if eos.sols == nil {
		eos.sols = new(Rank[T])
	}

	_ = eos.sols.Add(rank, elem)

	return nil
}

// AddErr adds an error to the EOS.
//
// Parameters:
//   - rank: The level of the error.
//   - err: The error to add.
//
// Returns:
//   - error: An error if the receiver is nil.
//
// Behaviors:
//   - If the error is nil, it is ignored.
//   - If at least a solution has been added, the error is ignored.
func (eos *ErrRorSol[T]) AddErr(rank int, err error) error {
	if err == nil {
		return nil
	} else if eos == nil {
		return common.ErrNilReceiver
	}

	if eos.errs == nil {
		return nil
	}

	_ = eos.errs.Add(rank, err)

	return nil
}

// HasError checks whether the EOS has an error.
//
// Returns:
//   - bool: True if the EOS has an error, false otherwise.
func (eos ErrRorSol[T]) HasError() bool {
	return eos.errs != nil && !eos.errs.IsEmpty()
}

// Errors returns the list of errors in descending order of rank.
//
// Returns:
//   - []error: The list of errors. Nil if there are no errors.
func (eos ErrRorSol[T]) Errors() []error {
	if eos.errs == nil {
		return nil
	}

	return eos.errs.Build()
}

// Sols returns the list of solutions in descending order of rank.
//
// Returns:
//   - []T: The list of solutions. Nil if there are no solutions.
func (eos ErrRorSol[T]) Sols() []T {
	return eos.sols.BuildAll()
}

// ChangeOrder changes the order in which elements are returned when methods such as
// `Sols` or `Errors` are called.
//
// Parameters:
//   - is_ascending: Whether to return elements in descending or ascending order.
//     If true, elements are returned in ascending order, otherwise in descending order.
//
// Returns:
//   - error: An error if the receiver is nil.
func (eos *ErrRorSol[T]) ChangeOrder(is_ascending bool) error {
	if eos == nil {
		return common.ErrNilReceiver
	}

	eos.errs.ChangeOrder(is_ascending)
	eos.sols.ChangeOrder(is_ascending)

	return nil
}
