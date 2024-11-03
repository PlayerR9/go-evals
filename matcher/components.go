package matcher

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/PlayerR9/go-evals/common"
)

// matchSingle is a matcher that matches a single element.
type matchSingle[I comparable] struct {
	// target is the element to match.
	target I

	// is_done is a flag that indicates if the matcher is done.
	is_done bool
}

// Close implements the Matcher interface.
func (m *matchSingle[I]) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	} else if !m.is_done {
		return NewErrNotAsExpected(true, "", nil, m.target)
	} else {
		return nil
	}
}

// Match implements the Matcher interface.
func (m *matchSingle[I]) Match(elem I) error {
	if m == nil {
		return common.ErrNilReceiver
	} else if m.is_done {
		return ErrMatchDone
	}

	if elem != m.target {
		return NewErrNotAsExpected(true, "", &elem, m.target)
	}

	m.is_done = true

	return nil
}

// Reset implements the Matcher interface.
func (m *matchSingle[I]) Reset() {
	if m == nil {
		return
	}

	m.is_done = false
}

// Matched implements the Matcher interface.
func (m matchSingle[I]) Matched() []I {
	if !m.is_done {
		return nil
	}

	return []I{m.target}
}

// Single returns a new matcher that matches a single element.
//
// Parameters:
//   - elem: The element to match.
//
// Returns:
//   - Matcher: The matcher. Never returns nil.
func Single[I comparable](elem I) Matcher[I] {
	return &matchSingle[I]{
		target:  elem,
		is_done: false,
	}
}

// matchSlice is a matcher that matches a literal.
type matchSlice[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// chars are the elements to match.
	chars []I

	// pos is the current position in the chars slice.
	pos int
}

// Close implements the Matcher interface.
func (w *matchSlice[I]) Close() error {
	if w == nil {
		return common.ErrNilReceiver
	} else if w.pos < len(w.chars) {
		return NewErrNotAsExpected(true, "char", nil, w.chars[w.pos])
	} else {
		return nil
	}
}

// Match implements the Matcher interface.
func (w *matchSlice[I]) Match(char I) error {
	if w == nil {
		return common.ErrNilReceiver
	} else if w.pos >= len(w.chars) {
		return ErrMatchDone
	}

	ok := w.chars[w.pos] == char
	if !ok {
		return NewErrNotAsExpected(true, "char", &char, w.chars[w.pos])
	}

	w.pos++

	w.matched = append(w.matched, char)

	return nil
}

// Reset implements the Matcher interface.
func (w *matchSlice[I]) Reset() {
	if w == nil {
		return
	}

	w.pos = 0

	if len(w.matched) > 0 {
		clear(w.matched)
		w.matched = nil
	}
}

// Matched implements the Matcher interface.
func (m matchSlice[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Slice returns a Matcher that matches a given slice. If the slice
// has only one element, then the element is matched directly.
//
// Parameters:
//   - slice: The slice to match.
//
// Returns:
//   - Matcher: The matcher. Nil if the word is empty.
func Slice[I comparable](slice []I) Matcher[I] {
	if len(slice) == 0 {
		return nil
	} else if len(slice) == 1 {
		return Single(slice[0])
	}

	return &matchSlice[I]{
		chars: slice,
		pos:   0,
	}
}

// matchFn is a matcher that matches a group of elements.
type matchFn[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// groupFn is the function to match.
	groupFn Predicate[I]

	// group_name is the name of the group.
	group_name string

	// is_done is a flag that indicates if the matcher is done.
	is_done bool
}

// Close implements the Matcher interface.
func (m *matchFn[I]) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	} else if !m.is_done {
		return fmt.Errorf("expected %s, got nothing", m.group_name)
	} else {
		return nil
	}
}

// Match implements the Matcher interface.
func (m *matchFn[I]) Match(elem I) error {
	if m == nil {
		return common.ErrNilReceiver
	} else if m.is_done {
		return ErrMatchDone
	}

	if !m.groupFn(elem) {
		return fmt.Errorf("expected %s, got %v", m.group_name, elem)
	}

	m.is_done = true

	m.matched = append(m.matched, elem)

	return nil
}

// Reset implements the Matcher interface.
func (m *matchFn[I]) Reset() {
	if m == nil {
		return
	}

	m.is_done = false

	if len(m.matched) > 0 {
		clear(m.matched)
		m.matched = nil
	}
}

// Matched implements the Matcher interface.
func (m matchFn[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Fn returns a new matcher that matches according to a predicate.
//
// Parameters:
//   - group_name: The name of the group.
//   - predicate: The function to match.
//
// Returns:
//   - Matcher: The matcher. Nil if the predicate is nil.
func Fn[I comparable](group_name string, predicate Predicate[I]) Matcher[I] {
	if predicate == nil {
		return nil
	}

	return &matchFn[I]{
		group_name: group_name,
		groupFn:    predicate,
		is_done:    false,
	}
}

// Group returns a new matcher that matches a group of elements.
//
// Parameters:
//   - elems: The elements to match.
//
// Returns:
//   - Matcher: The matcher. If no elements are provided, then nil is
//     returned.
//
// If only one character is provided, then the character is matched directly.
func Group[I comparable](group_name string, elems []I) Matcher[I] {
	if len(elems) == 0 {
		return nil
	} else if len(elems) == 1 {
		return Single(elems[0])
	}

	fn := func(elem I) bool {
		ok := slices.Contains(elems, elem)
		return ok
	}

	return &matchFn[I]{
		group_name: group_name,
		groupFn:    fn,
		is_done:    false,
	}
}

// SortedGroup returns a new matcher that matches a group of elements.
//
// Parameters:
//   - elems: The elements to match.
//
// Returns:
//   - Matcher: The matcher. If no elements are provided, then nil is
//     returned.
//
// If only one character is provided, then the character is matched directly.
func SortedGroup[I cmp.Ordered](group_name string, elems []I) Matcher[I] {
	if len(elems) == 0 {
		return nil
	}

	elems = Sort(elems)

	if len(elems) == 1 {
		return Single(elems[0])
	}

	fn := func(elem I) bool {
		_, ok := slices.BinarySearch(elems, elem)
		return ok
	}

	return &matchFn[I]{
		group_name: group_name,
		groupFn:    fn,
		is_done:    false,
	}
}

// matchGreedy is a matcher that matches a given inner Matcher as many times as
// possible.
type matchGreedy[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// inner is the inner Matcher.
	inner Matcher[I]
}

// Close implements the Matcher interface.
func (w *matchGreedy[I]) Close() error {
	if w == nil {
		return common.ErrNilReceiver
	}

	err := w.inner.Close()
	if err != nil {
		return fmt.Errorf("while matching greedy: %w", err)
	}

	w.matched = append(w.matched, w.inner.Matched()...)

	return nil
}

// Match implements the Matcher interface.
func (w *matchGreedy[I]) Match(char I) error {
	if w == nil {
		return common.ErrNilReceiver
	}

	err := w.inner.Match(char)
	if err == nil {
		return nil
	} else if err != ErrMatchDone {
		return fmt.Errorf("while matching many: %w", err)
	}

	w.matched = append(w.matched, w.inner.Matched()...)
	w.inner.Reset()

	err = w.inner.Match(char)
	if err == nil {
		return nil
	} else if err == ErrMatchDone {
		panic(errors.New("inner Matcher should not return ErrMatchDone as its first match"))
	} else {
		return ErrMatchDone
	}
}

// Reset implements the Matcher interface.
func (w *matchGreedy[I]) Reset() {
	if w == nil {
		return
	}

	if len(w.matched) > 0 {
		clear(w.matched)
		w.matched = nil
	}

	w.inner.Reset()
}

// Matched implements the Matcher interface.
func (m matchGreedy[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Greedy returns a Matcher that matches a given inner Matcher as many times as
// possible until it either fails, returns ErrMatchDone, or the input stream is exhausted.
//
// Parameters:
//   - inner: The inner Matcher.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
func Greedy[I comparable](inner Matcher[I]) Matcher[I] {
	if inner == nil {
		return nil
	}

	return &matchGreedy[I]{
		inner: inner,
	}
}

// matchSequence represents a sequence of Matchers.
type matchSequence[I comparable] struct {
	// matched are the runes that have been matched so far.
	matched []I

	// seq is the sequence of Matchers.
	seq []Matcher[I]

	// idx is the index of the current Matcher in the sequence.
	idx int
}

// Close implements the Matcher interface.
func (w *matchSequence[I]) Close() error {
	if w == nil {
		return common.ErrNilReceiver
	}

	if w.idx >= len(w.seq) {
		return nil
	}

	m := w.seq[w.idx]

	err := m.Close()
	if err != nil {
		return fmt.Errorf("while matching sequence: %w", err)
	}

	w.matched = append(w.matched, m.Matched()...)

	if w.idx+1 < len(w.seq) {
		return errors.New("matching sequence is not complete")
	}

	w.idx++

	return nil
}

// Match implements the Matcher interface.
func (w *matchSequence[I]) Match(char I) error {
	if w == nil {
		return common.ErrNilReceiver
	}

	for w.idx < len(w.seq) {
		m := w.seq[w.idx]

		err := m.Match(char)
		if err == nil {
			return nil
		} else if err != ErrMatchDone {
			return err
		}

		w.matched = append(w.matched, m.Matched()...)

		w.idx++
	}

	return ErrMatchDone
}

// Reset implements the Matcher interface.
func (w *matchSequence[I]) Reset() {
	if w == nil {
		return
	}

	if len(w.matched) > 0 {
		clear(w.matched)
		w.matched = nil
	}

	w.idx = 0

	for _, m := range w.seq {
		m.Reset()
	}
}

// Matched implements the Matcher interface.
func (m matchSequence[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Sequence returns a Matcher that matches a sequence of provided Matchers
// in the order they are given. The sequence will be processed by iterating
// through each Matcher and attempting to match the input element.
//
// Parameters:
//   - seq: A variadic number of Matcher instances. Matchers in the sequence are
//     expected to be non-nil objects.
//
// Returns:
//   - Matcher: A Matcher that represents a sequence of Matchers. Returns nil
//     if no non-nil Matchers are provided.
//
// If only one non-nil Matcher is provided, that Matcher is returned as it.
//
// A sequence can arbitrarily stop at any matcher in the sequence only if
// it is valid at that point.
func Sequence[I comparable](seq ...Matcher[I]) Matcher[I] {
	_ = RejectNils(&seq)
	if len(seq) == 0 {
		return nil
	} else if len(seq) == 1 {
		return seq[0]
	}

	return &matchSequence[I]{
		seq: seq,
		idx: 0,
	}
}

// matchWithBound is a matcher that matches a given inner Matcher until a
// boundary element is encountered.
type matchWithBound[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// match_inner is a flag that indicates that the inner Matcher should be
	// matched again.
	match_inner bool

	// inner is the inner Matcher.
	inner Matcher[I]

	// bound is the boundary function.
	bound Predicate[I]
}

// Close implements the Matcher interface.
func (m *matchWithBound[I]) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if !m.match_inner {
		return nil
	}

	err := m.inner.Close()
	if err != nil {
		return fmt.Errorf("while matching with boundary: %w", err)
	}

	m.match_inner = false

	m.matched = append(m.matched, m.inner.Matched()...)

	return nil
}

// Match implements the Matcher interface.
func (m *matchWithBound[I]) Match(char I) error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if m.match_inner {
		err := m.inner.Match(char)
		if err == nil {
			return nil
		} else if err != ErrMatchDone {
			return fmt.Errorf("while matching with boundary: %w", err)
		}

		m.match_inner = false

		m.matched = append(m.matched, m.inner.Matched()...)
	}

	ok := m.bound(char)
	if ok {
		return ErrMatchDone
	}

	return errors.New("boundary not satisfied")
}

// Reset implements the Matcher interface.
func (m *matchWithBound[I]) Reset() {
	if m == nil {
		return
	}

	if len(m.matched) > 0 {
		clear(m.matched)
		m.matched = nil
	}

	m.inner.Reset()

	m.match_inner = true
}

// Matched implements the Matcher interface.
func (m matchWithBound[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// WithBound returns a Matcher that matches a given inner Matcher until a
// boundary element is encountered.
//
// Parameters:
//   - inner: The inner Matcher.
//   - bound: The boundary function that defines the boundary that needs to be
//     satisfied for the Matcher to return ErrMatchDone.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
//
// If the bound function is nil, the inner Matcher is returned as is.
func WithBound[I comparable](inner Matcher[I], bound Predicate[I]) Matcher[I] {
	if inner == nil {
		return nil
	} else if bound == nil {
		return inner
	}

	return &matchWithBound[I]{
		inner: inner,
		bound: bound,

		match_inner: true,
	}
}

// matchAutoBound is a matcher that matches a given inner Matcher if and only if
// the next element does not satisfy the inner Matcher.
type matchAutoBound[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// match_inner is a flag that indicates that the inner Matcher should be
	// matched again.
	match_inner bool

	// inner is the inner Matcher.
	inner Matcher[I]
}

// Close implements the Matcher interface.
func (m *matchAutoBound[I]) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if !m.match_inner {
		return nil
	}

	err := m.inner.Close()
	if err != nil {
		return fmt.Errorf("while matching with boundary: %w", err)
	}

	m.match_inner = false

	m.matched = append(m.matched, m.inner.Matched()...)
	m.inner.Reset()

	return nil
}

// Match implements the Matcher interface.
func (m *matchAutoBound[I]) Match(elem I) error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if m.match_inner {
		err := m.inner.Match(elem)
		if err == nil {
			return nil
		} else if err != ErrMatchDone {
			return fmt.Errorf("while matching with boundary: %w", err)
		}

		m.match_inner = false

		m.matched = append(m.matched, m.inner.Matched()...)

		m.inner.Reset()
	}

	err := m.inner.Match(elem)
	if err != nil {
		return ErrMatchDone
	}

	return errors.New("boundary not satisfied")
}

// Reset implements the Matcher interface.
func (m *matchAutoBound[I]) Reset() {
	if m == nil {
		return
	}

	if len(m.matched) > 0 {
		clear(m.matched)
		m.matched = nil
	}

	m.inner.Reset()

	m.match_inner = true
}

// Matched implements the Matcher interface.
func (m matchAutoBound[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// AutoBound returns a Matcher that matches a given inner Matcher if and only if
// the next element does not satisfy the inner Matcher. This is useful when
// disambiguating between prefixes rules.
//
// Parameters:
//   - inner: The inner Matcher.
//
// Returns:
//   - Matcher: The matcher. Nil if the inner matcher is nil.
//
// It is equivalent to:
//
//	WithRightBound(inner, func(char element) bool { err := inner.Match(char); return err != nil }).
func AutoBound[I comparable](inner Matcher[I]) Matcher[I] {
	if inner == nil {
		return nil
	}

	return &matchAutoBound[I]{
		inner:       inner,
		match_inner: true,
	}
}

// Range returns a matcher that matches a group of elements between left and
// right.
//
// If left and right are equal, the returned matcher will match exactly one
// element. Otherwise, the returned matcher will match any element in the
// range [left, right].
//
// Parameters:
//   - left: The left boundary of the group.
//   - right: The right boundary of the group.
//
// Returns:
//   - Matcher: The matcher. Never returns nil.
func Range[I cmp.Ordered](left, right I) Matcher[I] {
	if left == right {
		return Single(left)
	}

	if left > right {
		left, right = right, left
	}

	return Fn(
		fmt.Sprintf("[%v-%v]", left, right),
		func(elem I) bool {
			return elem >= left && elem <= right
		},
	)
}

// matchOr is a matcher that matches any of the given matchers; prioritizing
// the longest valid match and the first one added when there is a tie.
type matchOr[I comparable] struct {
	// matched are the elements that were matched.
	matched []I

	// matchers are the matchers to be matched.
	matchers []Matcher[I]

	// indices are the indices of the matchers that were matched.
	indices []int

	// sol is the index of the solution. -1 if there is no solution.
	sol int

	// errs are the most recent errors.
	errs []error
}

// Close implements the Matcher interface.
func (m *matchOr[I]) Close() error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if len(m.indices) == 0 {
		if m.sol == -1 {
			return errors.Join(m.errs...)
		} else {
			return nil
		}
	}

	var errs []error
	var top int

	for _, idx := range m.indices {
		match := m.matchers[idx]

		err := match.Close()
		if err != nil {
			errs = append(errs, err)
		} else {
			m.matched = match.Matched()
			m.sol = idx
		}
	}

	m.indices = m.indices[:top:top]

	if len(errs) > 0 {
		m.errs = errs
	}

	return nil
}

// Match implements the Matcher interface.
func (m *matchOr[I]) Match(elem I) error {
	if m == nil {
		return common.ErrNilReceiver
	}

	if len(m.indices) == 0 {
		if m.sol == -1 {
			return errors.Join(m.errs...)
		} else {
			return ErrMatchDone
		}
	}

	var errs []error
	var top int

	for _, idx := range m.indices {
		match := m.matchers[idx]

		err := match.Match(elem)
		if err == nil {
			m.indices[top] = idx
			top++
		} else if err == ErrMatchDone {
			m.matched = match.Matched()
			m.sol = idx
		} else if m.sol != -1 {
			errs = append(errs, err)
		}
	}

	m.indices = m.indices[:top:top]

	if len(errs) > 0 {
		m.errs = errs
	}

	return nil
}

// Reset implements the Matcher interface.
func (m *matchOr[I]) Reset() {
	if m == nil {
		return
	}

	if len(m.indices) > 0 {
		clear(m.indices)
	}

	m.indices = make([]int, 0, len(m.matchers))
	for i := range m.matchers {
		m.indices = append(m.indices, i)
	}
	slices.Reverse(m.indices)

	for _, match := range m.matchers {
		match.Reset()
	}

	m.sol = -1

	if len(m.errs) > 0 {
		clear(m.errs)
		m.errs = nil
	}

	if len(m.matched) > 0 {
		clear(m.matched)
		m.matched = nil
	}
}

// Matched implements the Matcher interface.
func (m matchOr[I]) Matched() []I {
	if len(m.matched) == 0 {
		return nil
	}

	matched := make([]I, len(m.matched))
	copy(matched, m.matched)

	return matched
}

// Or returns a Matcher that matches any of the given matchers; returning the longest
// valid match.
//
// The matchers are given in descending order of priority, with the first
// matcher being the highest priority. If multiple matchers match, the
// highest priority matcher is considered to have matched.
//
// Parameters:
//   - matchers: The matchers to be matched. Nil elements are discarded in-place.
//
// Returns:
//   - Matcher: The matcher. Return nil if no matchers are given or are all nil.
//
// If a single non-nil matcher is given, the returned Matcher is the given matcher itself
// as is since there is no need to wrap it in an Or matcher.
func Or[I comparable](matchers ...Matcher[I]) Matcher[I] {
	_ = RejectNils(&matchers)
	if len(matchers) == 0 {
		return nil
	} else if len(matchers) == 1 {
		return matchers[0]
	}

	m := &matchOr[I]{
		matchers: matchers,
		sol:      -1,
	}

	m.indices = make([]int, 0, len(matchers))
	for i := range matchers {
		m.indices = append(m.indices, i)
	}

	slices.Reverse(m.indices)

	return m
}
