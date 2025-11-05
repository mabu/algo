// Package sorted provides data structures that maintain the elements in a sorted order.
package sorted

import (
	"cmp"
	"iter"
)

// SetItem refers to an element in the Set.
type SetItem[T any] struct {
	l, r, parent *SetItem[T]
	level        int8
	value        T
	set          *Set[T]
}

// Value returns the item's value.
func (si *SetItem[T]) Value() T {
	return si.value
}

// Next returns the next larger item in the set and true,
// or nil and false if si is already the largest element.
func (si *SetItem[T]) Next() (*SetItem[T], bool) {
	if si.r.level == 0 {
		for si.parent != nil && si == si.parent.r {
			si = si.parent
		}
		return si.parent, si.parent != nil
	}
	si = si.r
	for si.l.level > 0 {
		si = si.l
	}
	return si, true
}

// Prev returns the next smaller item in the set and true,
// or nil and false if si is already the smallest element.
func (si *SetItem[T]) Prev() (*SetItem[T], bool) {
	if si.l.level == 0 {
		for si.parent != nil && si == si.parent.l {
			si = si.parent
		}
		return si.parent, si.parent != nil
	}
	si = si.l
	for si.r.level > 0 {
		si = si.r
	}
	return si, true
}

// InsertNearby inserts x to the same set as si.
// It is more efficient than calling [Set.Insert]
// if x would end up right before or after si.
//
// Returns the SetItem whose value is x,
// and a bool indicating whether this is a newly added item.
func (si *SetItem[T]) InsertNearby(x T) (*SetItem[T], bool) {
	return si.set.insertNearby(si, x)
}

// Set is a sorted (ordered) set of T.
type Set[T any] struct {
	root   *SetItem[T]
	bottom *SetItem[T]
	finder[T]
	size int
}

func (s *Set[T]) skew(t *SetItem[T]) *SetItem[T] {
	if t.l.level == t.level {
		if p := t.parent; p == nil {
			s.root = t.l
		} else {
			if p.l == t {
				p.l = t.l
			} else {
				p.r = t.l
			}
		}
		t.l, t.parent, t.l.r, t.l.parent, t.l.r.parent, t = t.l.r, t.l, t, t.parent, t, t.l
	}
	return t
}

func (s *Set[T]) split(t *SetItem[T]) (*SetItem[T], bool) {
	if t.level == t.r.r.level {
		if p := t.parent; p == nil {
			s.root = t.r
		} else {
			if p.l == t {
				p.l = t.r
			} else {
				p.r = t.r
			}
		}
		t.r, t.parent, t.r.l, t.r.parent, t.r.l.parent, t = t.r.l, t.r, t, t.parent, t, t.r
		t.level++
		return t, true
	}
	return t, false
}

func newSet[T any](f finder[T]) *Set[T] {
	bottom := new(SetItem[T])
	bottom.l = bottom
	bottom.r = bottom
	return &Set[T]{
		root:   bottom,
		bottom: bottom,
		finder: f,
	}
}

// finder contains the methods that depend on the comparator.
// We could as well have only the comparator itself as a func in the Set,
// but this appears to be more efficient.
type finder[T any] interface {
	// find looks for a value x in the set. If found, returns that node and nil.
	// Otherwise returns the node that was the last on the path towards
	// where x would be found, and the pointer to its l or r field depending on
	// how x compares. That is, if target != nil, then the new node can be
	// added to *target, and last used as its parent.
	// The given root will never be bottom.
	find(root *SetItem[T], x T) (last *SetItem[T], target **SetItem[T])

	findGreaterThanOrEqual(root *SetItem[T], x T) (*SetItem[T], bool)

	insertNearby(*SetItem[T], T) (*SetItem[T], bool)
}

// NewSet creates a new sorted set of T, using < for comparisons.
func NewSet[T cmp.Ordered]() *Set[T] {
	return newSet[T](set[T]{})
}

// NewSetFunc creates a new set of T which is ordered according to cmp.
//
// If T is or contains a pointer,
// the values referenced by it must not be changed
// in a way that affects the order
// for as long as it is in the Set.
func NewSetFunc[T any](cmp func(T, T) int) *Set[T] {
	return newSet[T](setFunc[T](cmp))
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return s.size
}

// First returns the smallest SetItem and true, or nil and false if the set is empty.
func (s *Set[T]) First() (*SetItem[T], bool) {
	if s.root.level == 0 {
		return nil, false
	}
	v := s.root
	for v.l.level > 0 {
		v = v.l
	}
	return v, true
}

// Last returns the largest SetItem and true, or nil and false if the set is empty.
func (s *Set[T]) Last() (*SetItem[T], bool) {
	if s.root.level == 0 {
		return nil, false
	}
	v := s.root
	for v.r.level > 0 {
		v = v.r
	}
	return v, true
}

// Min returns the smallest value in the set and true.
// If the set is empty returns false.
func (s *Set[T]) Min() (T, bool) {
	si, ok := s.First()
	if !ok {
		var v T
		return v, false
	}
	return si.value, true
}

// Max returns the largest value in the set and true.
// If the set is empty returns false.
func (s *Set[T]) Max() (T, bool) {
	si, ok := s.Last()
	if !ok {
		var v T
		return v, false
	}
	return si.value, true
}

// All returns an iterator over all elements in the set in sorted order.
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		si, ok := s.First()
		for ok && yield(si.value) {
			si, ok = si.Next()
		}
	}
}

// Backward returns an iterator over all elements in the set in reverse order.
func (s *Set[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		si, ok := s.Last()
		for ok && yield(si.value) {
			si, ok = si.Prev()
		}
	}
}

// Insert adds x to the set.
// Returns whether the insertion happened, i.e.
// returns false if an element that compares as equal to x was already in the set, otherwise returns true.
func (s *Set[T]) Insert(x T) (added bool) {
	if s.root == s.bottom {
		s.root = s.newItem(x, nil)
		return true
	}
	_, added = s.insert(s.root, x)
	return added
}

// root != bottom
func (s *Set[T]) insert(root *SetItem[T], x T) (*SetItem[T], bool) {
	last, target := s.find(s.root, x)
	if target == nil {
		return last, false
	}
	result := s.newItem(x, last)
	*target = result
	for t, n := last, 0; t != nil && n < 2; t = t.parent {
		t = s.skew(t)
		var ok bool
		t, ok = s.split(t)
		if ok {
			n = 0
		} else {
			n++
		}
		if t.parent == nil {
			s.root = t
		}
	}
	return result, true
}

func (s *Set[T]) newItem(x T, parent *SetItem[T]) *SetItem[T] {
	s.size++
	return &SetItem[T]{
		value:  x,
		level:  1,
		l:      s.bottom,
		r:      s.bottom,
		parent: parent,
		set:    s,
	}
}

// Has reports whether x is in the set.
func (s *Set[T]) Has(x T) bool {
	if s.root == s.bottom {
		return false
	}
	_, target := s.find(s.root, x)
	return target == nil
}

// FindGreaterThanOrEqual returns the first SetItem that is greater than or equal to x, and true.
// If there is no such element, returns false.
func (s *Set[T]) FindGreaterThanOrEqual(x T) (*SetItem[T], bool) {
	if s.root == s.bottom {
		return nil, false
	}
	return s.findGreaterThanOrEqual(s.root, x)
}

// Delete removes x from the set if it exists.
// The return value indicates whether the removal happened.
func (s *Set[T]) Delete(x T) (deleted bool) {
	defer func() {
		if deleted {
			s.size--
		}
	}()
	if s.root == s.bottom {
		return false
	}
	last, target := s.find(s.root, x)
	if target != nil {
		return false
	}
	if last.level > 1 {
		successor := last.r
		for successor.l != s.bottom {
			successor = successor.l
		}
		last.value = successor.value
		last = successor
	} else if last.parent == nil {
		s.root = last.r
		s.root.parent = nil
		return true
	}
	// Level 1, not root.
	last.r.parent = last.parent
	if last.parent.l == last {
		last.parent.l = last.r
	} else {
		last.parent.r = last.r
	}
	last = last.parent
	for {
		var ok bool
		last, ok = s.decreaseLevel(last)
		if !ok {
			break
		}
		if last.parent == nil {
			s.root = last
			break
		}
		last = last.parent
	}
	return true
}

func (s *Set[T]) decreaseLevel(t *SetItem[T]) (*SetItem[T], bool) {
	if t.level > t.l.level+1 || t.level > t.r.level+1 {
		t.level--
		if t.r.level > t.level {
			t.r.level = t.level
		}

		t = s.skew(t)
		t.r = s.skew(t.r)
		t.r.r = s.skew(t.r.r)
		t, _ = s.split(t)
		if t.r.level != 0 {
			t.r, _ = s.split(t.r)
		}
		return t, true
	}
	return t, false
}

type (
	set[T cmp.Ordered] struct{}
	setFunc[T any]     func(T, T) int
)

func (set[T]) find(root *SetItem[T], x T) (last *SetItem[T], target **SetItem[T]) {
	for root.value != x {
		if x < root.value {
			if root.l.level == 0 {
				return root, &root.l
			}
			root = root.l
		} else {
			if root.r.level == 0 {
				return root, &root.r
			}
			root = root.r
		}
	}
	return root, nil
}

func (cmp setFunc[T]) find(root *SetItem[T], x T) (last *SetItem[T], target **SetItem[T]) {
	for {
		switch c := cmp(x, root.value); {
		case c < 0:
			if root.l.level == 0 {
				return root, &root.l
			}
			root = root.l
		case c > 0:
			if root.r.level == 0 {
				return root, &root.r
			}
			root = root.r
		default:
			return root, nil
		}
	}
}

func (set[T]) findGreaterThanOrEqual(root *SetItem[T], x T) (*SetItem[T], bool) {
	for root.value != x {
		if x < root.value {
			if root.l.level == 0 {
				return root, true
			}
			root = root.l
		} else {
			if root.r.level == 0 {
				return root.Next()
			}
			root = root.r
		}
	}
	return root, true
}

func (cmp setFunc[T]) findGreaterThanOrEqual(root *SetItem[T], x T) (*SetItem[T], bool) {
	for {
		switch c := cmp(x, root.value); {
		case c < 0:
			if root.l.level == 0 {
				return root, true
			}
			root = root.l
		case c > 0:
			if root.r.level == 0 {
				return root.Next()
			}
			root = root.r
		default:
			return root, true
		}
	}
}

func (set[T]) insertNearby(si *SetItem[T], x T) (*SetItem[T], bool) {
	if si.value == x {
		return si, false
	}
	if si.value < x {
		next, ok := si.Next()
		if !ok {
			return si.set.insert(si, x)
		}
		if next.value == x {
			return next, false
		}
		if next.value > x {
			return si.set.insert(lower(si, next), x)
		}
	} else {
		prev, ok := si.Prev()
		if !ok {
			return si.set.insert(si, x)
		}
		if prev.value == x {
			return prev, false
		}
		if prev.value < x {
			return si.set.insert(lower(si, prev), x)
		}
	}
	return si.set.insert(si.set.root, x)
}

func (cmp setFunc[T]) insertNearby(si *SetItem[T], x T) (*SetItem[T], bool) {
	siCmp := cmp(si.value, x)
	if siCmp == 0 {
		return si, false
	}
	if siCmp < 0 {
		next, ok := si.Next()
		if !ok {
			return si.set.insert(si, x)
		}
		nextCmp := cmp(next.value, x)
		if nextCmp == 0 {
			return next, false
		}
		if nextCmp > 0 {
			return si.set.insert(lower(si, next), x)
		}
	} else {
		prev, ok := si.Prev()
		if !ok {
			return si.set.insert(si, x)
		}
		prevCmp := cmp(prev.value, x)
		if prevCmp == 0 {
			return prev, false
		}
		if prevCmp < 0 {
			return si.set.insert(lower(si, prev), x)
		}
	}
	return si.set.insert(si.set.root, x)
}

func lower[T any](a, b *SetItem[T]) *SetItem[T] {
	if a.level < b.level {
		return a
	}
	if b.level < a.level {
		return b
	}
	if a.parent == b {
		return b
	}
	return a
}
