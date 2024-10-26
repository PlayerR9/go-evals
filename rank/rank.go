package rank

import (
	"iter"
	"slices"

	"github.com/PlayerR9/go-evals/common"
)

// Rank is a data structure that holds a list of solutions in order of rank.
// An empty rank can either be created with the `var rank Rank[T]` syntax
// or with the `new(Rank[T])` constructor.
type Rank[T any] struct {
	// buckets maps ranks to buckets.
	buckets map[int][]T

	// ranks is a sorted list of ranks.
	ranks []int

	// size is the number of elements in the rank.
	size int

	// is_ascending is a flag that indicates whether the elements should
	// be returned in descending or ascending order when methods such as
	// `Elem` or `Build` are called.
	is_ascending bool
}

// Size returns the number of elements in the rank.
//
// Returns:
//   - int: The number of elements in the rank. Never negative.
func (r Rank[T]) Size() int {
	return r.size
}

// IsEmpty checks whether the rank is empty.
//
// Returns:
//   - bool: True if the rank is empty, false otherwise.
func (r Rank[T]) IsEmpty() bool {
	return r.size == 0
}

// Reset resets the rank for reuse.
func (r *Rank[T]) Reset() {
	if r == nil {
		return
	}

	if len(r.buckets) > 0 {
		for k, v := range r.buckets {
			clear(v)
			r.buckets[k] = nil
		}

		clear(r.buckets)
		r.buckets = nil
	}

	if len(r.ranks) > 0 {
		clear(r.ranks)
		r.ranks = nil
	}

	r.size = 0
}

// Add adds an element to the rank.
//
// Parameters:
//   - rank: The rank of the element.
//   - elem: The element to add.
//
// Returns:
//   - error: An error if the receiver is nil.
func (r *Rank[T]) Add(rank int, elem T) error {
	if r == nil {
		return common.ErrNilReceiver
	}

	if r.buckets == nil {
		r.buckets = make(map[int][]T)
	}

	bucket, ok := r.buckets[rank]
	if !ok {
		r.buckets[rank] = []T{elem}

		pos, ok := slices.BinarySearch(r.ranks, rank)
		if !ok {
			r.ranks = slices.Insert(r.ranks, pos, rank)
		}
	} else {
		r.buckets[rank] = append(bucket, elem)
	}

	r.size++

	return nil
}

// Elem iterates over the elements in the rank in rank order. The order of elements
// within a rank is guaranteed to be the same as the order in which they were
// added. However, the exact order of elements with different ranks is in descending
// order unless the method `ChangeOrder(true)` is called.
//
// Returns:
//   - iter.Seq2[int, T]: The elements in the rank. Never returns nil.
func (r Rank[T]) Elem() iter.Seq2[int, T] {
	if r.size == 0 {
		return func(yield func(int, T) bool) {}
	}

	var fn iter.Seq2[int, T]

	if r.is_ascending {
		fn = func(yield func(int, T) bool) {
			for _, rank := range r.ranks {
				bucket := r.buckets[rank]

				for _, elem := range bucket {
					if !yield(rank, elem) {
						return
					}
				}
			}
		}
	} else {
		fn = func(yield func(int, T) bool) {
			for i := len(r.ranks) - 1; i >= 0; i-- {
				rank := r.ranks[i]

				bucket := r.buckets[rank]

				for _, elem := range bucket {
					if !yield(rank, elem) {
						return
					}
				}
			}
		}
	}

	return fn
}

// ChangeOrder changes the order in which elements are returned when methods such as
// `Elem` or `Build` are called.
//
// Parameters:
//   - is_ascending: Whether to return elements in descending or ascending order.
//     If true, elements are returned in ascending order, otherwise in descending order.
//
// Returns:
//   - error: An error if the receiver is nil.
func (r *Rank[T]) ChangeOrder(is_ascending bool) error {
	if r == nil {
		return common.ErrNilReceiver
	}

	r.is_ascending = is_ascending

	return nil
}

// Build is a more efficent way of calling slices.Collect(r.Elem()).
//
// Returns:
//   - []T: The slice of elements in the rank. Returns nil if the rank is empty.
func (r Rank[T]) Build() []T {
	if r.size == 0 {
		return nil
	}

	slice := make([]T, 0, r.size)

	if r.is_ascending {
		for _, rank := range r.ranks {
			bucket := r.buckets[rank]

			for i := len(bucket) - 1; i >= 0; i-- {
				slice = append(slice, bucket[i])
			}
		}
	} else {
		for i := len(r.ranks) - 1; i >= 0; i-- {
			rank := r.ranks[i]

			bucket := r.buckets[rank]

			for _, elem := range bucket {
				slice = append(slice, elem)
			}
		}
	}

	return slice
}
