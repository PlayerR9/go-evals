package common

// Close closes the given channel safely by checking if it's not nil and not closed yet.
// If the channel is already closed or nil, this function does nothing.
//
// This function is designed to be idempotent, meaning it can be safely called multiple times on the same channel.
// It's particularly useful for deferring the closing of a channel.
//
// Parameters:
//   - ch: The channel to close.
func Close[T any](ch chan T) {
	if ch == nil {
		return
	}

	select {
	case _, ok := <-ch:
		if !ok {
			return
		}
	default:
	}

	close(ch)
}

// Send sends the given element on the channel if the channel is not nil and not closed.
// If the channel is closed or the buffer is full, it does nothing.
//
// Parameters:
//   - ch: The channel to send the element on.
//   - elem: The element to send.
//
// Returns:
//   - bool: True if the element was sent, false otherwise.
func Send[T any](ch chan T, elem T) bool {
	if ch == nil {
		return false
	}

	select {
	case ch <- elem:
		return true
	default:
		return false
	}
}
