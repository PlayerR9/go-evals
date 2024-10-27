package catcher

import (
	"context"

	"github.com/PlayerR9/go-evals/common"
)

// listenChAct is an action that listens to a channel.
type listenChAct[T any] struct {
	// ch is the channel to listen to.
	ch <-chan T
}

// Run implements the contexts.Action interface.
func (act *listenChAct[T]) Run(ctx context.Context) error {
	if act == nil {
		return common.ErrNilReceiver
	}

	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for elem := range act.ch {
			_ = c.addElem(elem)
		}
	}()

	return nil
}

// ListenCh returns an action that listens to the given channel and adds all
// received elements to the catcher.
//
// Parameters:
//   - ch: The channel to listen to.
//
// Returns:
//   - contexts.Action: An action that listens to the channel. Nil if ch is nil.
func ListenCh[T any](ch <-chan T) Action {
	if ch == nil {
		return nil
	}

	return &listenChAct[T]{
		ch: ch,
	}
}

// sendElemAct is an action that sends an element on a channel.
type sendElemAct[T any] struct {
	// elem is the element to send.
	elem T
}

// Run implements the contexts.Action interface.
func (act *sendElemAct[T]) Run(ctx context.Context) error {
	if act == nil {
		return common.ErrNilReceiver
	}

	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	_ = c.addElem(act.elem)

	return nil
}

// SendElem returns an action that sends the given element on the channel.
//
// Parameters:
//   - elem: The element to send.
//
// Returns:
//   - contexts.Action: An action that sends the element. Never returns nil.
func SendElem[T any](elem T) Action {
	return &sendElemAct[T]{
		elem: elem,
	}
}

// getElemsAct is an action that gets all elements from the catcher.
type getElems[T any] struct {
	// dest is the destination for the elements.
	dest *[]T
}

// Run implements the contexts.Action interface.
func (act *getElems[T]) Run(ctx context.Context) error {
	if act == nil {
		return common.ErrNilReceiver
	}

	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	elems := c.getElems()
	*act.dest = elems

	return nil
}

// GetElems returns an action that gets all elements from the catcher and stores
// them in the destination.
//
// Parameters:
//   - dest: The destination for the elements.
//
// Returns:
//   - contexts.Action: An action that gets all elements from the catcher and stores
//     them in the destination. nil if dest is nil.
func GetElems[T any](dest *[]T) Action {
	if dest == nil {
		return nil
	}

	return &getElems[T]{
		dest: dest,
	}
}

// waitAct is an action that waits for the catcher to be done listening to the
// channel.
type waitAct[T any] struct{}

// Run implements the contexts.Action interface.
func (act *waitAct[T]) Run(ctx context.Context) error {
	if act == nil {
		return common.ErrNilReceiver
	}

	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	c.wg.Wait()

	return nil
}

// Wait returns an action that waits for the catcher to complete its operations.
//
// Returns:
//   - sf.Action: An action that waits for all operations in the catcher to finish. Never returns nil.
func Wait[T any]() Action {
	return &waitAct[T]{}
}

// resetCacheAct is an action that clears the cache.
type resetCacheAct[T any] struct{}

// Run implements the contexts.Action interface.
func (act *resetCacheAct[T]) Run(ctx context.Context) error {
	if act == nil {
		return common.ErrNilReceiver
	}

	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	c.clear()

	return nil
}

// ResetCache returns an action that clears the cache.
//
// Returns:
//   - sf.Action: An action that clears the cache. Never returns nil.
func ResetCache[T any]() Action {
	return &resetCacheAct[T]{}
}
