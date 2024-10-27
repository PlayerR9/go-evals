package catcher

import (
	"context"

	"github.com/PlayerR9/go-evals/common"
)

// Action is an interface for actions that can be executed in a context.
type Action interface {
	// Run executes the action within the provided context.
	//
	// Parameters:
	//   - ctx: The context in which the action is executed.
	//
	// Returns:
	//   - error: An error if the action fails during execution.
	Run(ctx context.Context) error
}

// RejectNilActions removes all nil actions from the provided slice of actions.
//
// Parameters:
//   - acts: A pointer to a slice of actions. The function modifies this slice
//     in-place to exclude any nil entries. If the slice is nil or empty,
//     the function does nothing.
//
// Behavior:
//   - If all elements are nil or no elements are specified, the slice is set to nil.
func RejectNilActions(acts *[]Action) {
	if acts == nil || len(*acts) == 0 {
		return
	}

	var top int

	for _, act := range *acts {
		if act != nil {
			(*acts)[top] = act
			top++
		}
	}

	if top == 0 {
		clear(*acts)
		*acts = nil

		return
	}

	clear((*acts)[top:])
	*acts = (*acts)[:top]
}

// Run runs the given actions in-order and returns the first error encountered.
//
// Parameters:
//   - ctx: The context to use for the actions. If the context is canceled, the function
//     returns the error associated with the context.
//   - acts: Zero or more actions to run. If no actions are specified, the function
//     returns nil.
//
// Returns:
//   - error: The first error encountered while running the actions. If no error is
//     encountered, the function returns nil.
func Run(ctx context.Context, acts ...Action) error {
	RejectNilActions(&acts)
	if len(acts) == 0 {
		return nil
	} else if ctx == nil {
		return common.NewErrNilParam("ctx")
	}

	done := ctx.Done()

	for _, act := range acts {
		select {
		case <-done:
			return ctx.Err()
		default:
			err := act.Run(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
