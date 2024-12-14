package internal

import "github.com/PlayerR9/go-evals/common"

type Queue[E any] struct {
	elems []E
}

func (q Queue[E]) IsEmpty() bool {
	return len(q.elems) == 0
}

func (l *Queue[E]) Enqueue(e E) error {
	if l == nil {
		return common.ErrNilReceiver
	}

	l.elems = append(l.elems, e)

	return nil
}

func (l *Queue[E]) Dequeue() (E, error) {
	if l == nil {
		return *new(E), common.ErrNilReceiver
	}

	if len(l.elems) == 0 {
		return *new(E), ErrEmptyQueue
	}

	e := l.elems[0]
	l.elems = l.elems[1:]

	return e, nil
}
