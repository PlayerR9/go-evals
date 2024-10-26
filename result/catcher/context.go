package catcher

import (
	"context"
	"sync"

	"github.com/PlayerR9/go-evals/common"
)

// contextKey is the context key for the context.
type contextKey struct{}

// fromContext is a helper function that returns the underlying context.
//
// Parameters:
//   - ctx: The context to get the underlying context from.
//
// Returns:
//   - *Context[T]: The underlying context. Returns nil if the context is not found.
//   - error: An error if the context is not found.
func fromContext[T any](ctx context.Context) (*Context[T], error) {
	if ctx == nil {
		return nil, common.NewErrNilParam("ctx")
	}

	c, ok := ctx.Value(contextKey{}).(*Context[T])
	if !ok {
		return nil, common.NewErrBadParam("ctx", "must contain a *Context[T]")
	}

	return c, nil
}

// Context is the context for the catcher.
type Context[T any] struct {
	// elems is the list of elements.
	elems []T

	// mu is the mutex for the list of elements.
	mu sync.Mutex

	// wg is the wait group for the list of elements.
	wg sync.WaitGroup
}

// NewContext creates a new context that includes a *Context[T] value.
//
// If the parent context is nil, a new background context is created.
// The new context is associated with a new *Context[T] and stored with a
// contextKey.
//
// Parameters:
//   - parent: The parent context to use as the base for the new context.
//
// Returns:
//   - context.Context: The new context with the *Context[T] value.
func NewContext[T any](parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	c := &Context[T]{}

	ctx := context.WithValue(parent, contextKey{}, c)

	return ctx
}

// addElem adds a new element to the context.
//
// Parameters:
//   - elem: The new element to add.
//
// Returns:
//   - error: An error if the context is nil.
func (c *Context[T]) addElem(elem T) error {
	if c == nil {
		return common.ErrNilReceiver
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.elems = append(c.elems, elem)

	return nil
}

// clear clears the elements in the context.
func (c *Context[T]) clear() {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.elems) == 0 {
		return
	}

	zero := *new(T)

	for i := range c.elems {
		c.elems[i] = zero
	}

	c.elems = c.elems[:0]
}

// getElems returns a copy of the elements stored in the context.
//
// This method locks the context's mutex to ensure thread safety while copying
// the elements. If the context is nil or the elements slice is nil, it returns
// nil.
//
// Returns:
//   - []T: A copy of the elements stored in the context.
func (c *Context[T]) getElems() []T {
	if c == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.elems == nil {
		return nil
	}

	elems := make([]T, len(c.elems))
	copy(elems, c.elems)

	return elems
}
