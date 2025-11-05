package sorted

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"testing"
	"testing/quick"

	gcmp "github.com/google/go-cmp/cmp"
)

func TestInsert(t *testing.T) {
	for name, set := range map[string]*Set[int]{
		"NewSet":     NewSet[int](),
		"NewSetFunc": NewSetFunc(cmp.Compare[int]),
	} {
		t.Run(name, func(t *testing.T) {
			operations := []struct {
				x    int
				want bool
			}{
				{1, true},
				{1, false},
				{2, true},
				{3, true},
				{2, false},
				{3, false},
			}
			for i, op := range operations {
				if got := set.Insert(op.x); got != op.want {
					t.Errorf("Operation #%d: Insert(%d) = %t, want %t", i, op.x, got, op.want)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			testCases := []struct {
				insert  []int
				delete  int
				want    bool
				wantAll []int
			}{
				{
					delete: 42,
				},
				{
					insert:  []int{1},
					delete:  2,
					wantAll: []int{1},
				},
				{
					insert: []int{1},
					delete: 1,
					want:   true,
				},
				{
					insert:  []int{1, 2},
					delete:  1,
					want:    true,
					wantAll: []int{2},
				},
				{
					insert:  []int{1, 2},
					delete:  2,
					want:    true,
					wantAll: []int{1},
				},
				{
					insert:  []int{1, 2},
					delete:  3,
					want:    false,
					wantAll: []int{1, 2},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  3,
					wantAll: []int{1, 2, 4, 6, 7, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  1,
					want:    true,
					wantAll: []int{2, 4, 6, 7, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  2,
					want:    true,
					wantAll: []int{1, 4, 6, 7, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  4,
					want:    true,
					wantAll: []int{1, 2, 6, 7, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  6,
					want:    true,
					wantAll: []int{1, 2, 4, 7, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  7,
					want:    true,
					wantAll: []int{1, 2, 4, 6, 8},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  8,
					want:    true,
					wantAll: []int{1, 2, 4, 6, 7},
				},
				{
					insert:  []int{2, 4, 6, 1, 7, 8},
					delete:  9,
					want:    false,
					wantAll: []int{1, 2, 4, 6, 7, 8},
				},
			}

			for _, tc := range testCases {
				s := newSet()
				for _, in := range tc.insert {
					s.Insert(in)
				}
				if got := s.Delete(tc.delete); got != tc.want {
					t.Errorf("After inserting %v, s.Delete(%d) = %t, want %t", tc.insert, tc.delete, got, tc.want)
				}
				if diff := gcmp.Diff(tc.wantAll, slices.Collect(s.All())); diff != "" {
					t.Errorf("After inserting %v and performing s.Delete(%d), All() diff (-want +got):\n%s", tc.insert, tc.delete, diff)
				}
				if got, want := s.Len(), len(tc.wantAll); got != want {
					t.Errorf("After inserting %v and performing s.Delete(%d): Len() = %d, want %d", tc.insert, tc.delete, got, want)
				}
			}
		})
	}
}

func TestFullToEmpty(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			for name, insert := range map[string]func(*Set[int]) func(v int) bool{
				"Insert": func(s *Set[int]) func(v int) bool {
					return s.Insert
				},
				"InsertNearby": func(s *Set[int]) func(v int) bool {
					// Most of the time these will not actually be nearby, but it should still work.
					var si *SetItem[int]
					return func(v int) bool {
						if si == nil {
							res := s.Insert(v)
							var ok bool
							si, ok = s.First()
							if !ok || si == nil {
								t.Fatalf("After s.Insert(%d), there is still no s.First() item.", v)
							}
							return res
						}
						var res bool
						si, res = si.InsertNearby(v)
						return res
					}
				},
			} {
				t.Run(name, func(t *testing.T) {
					s := newSet()
					insert := insert(s)
					if got := s.Len(); got != 0 {
						t.Errorf("newSet.Len() = %d, want 0", got)
					}
					for _, v := range permutation1 {
						if !insert(v) {
							t.Errorf("While inserting permutation: Insert(%d) = false, want true", v)
						}
					}
					if got, want := s.Len(), len(permutation1); got != want {
						t.Errorf("After inserting %d elements to a set, Len() = %d, want %d", want, got, want)
					}
					for _, v := range permutation2 {
						if !s.Delete(v) {
							t.Errorf("While deleting permutation: Delete(%d) = false, want true", v)
						}
					}
					if got := s.Len(); got != 0 {
						t.Errorf("After inserting many elements and deleting them all, Len() = %d, want 0", got)
					}
				})
			}
		})
	}
}

func TestHas(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			testCases := []struct {
				insert []int
				delete []int
				x      int
				want   bool
			}{
				{
					x: 10,
				},
				{
					insert: []int{1},
					x:      2,
				},
				{
					insert: []int{1, 2},
					x:      2,
					want:   true,
				},
				{
					insert: []int{1, 2},
					x:      1,
					want:   true,
				},
				{
					insert: []int{1, 2},
					x:      3,
				},
				{
					insert: []int{1, 2},
					x:      0,
				},
				{
					insert: []int{2, 3},
					x:      1,
				},
				{
					insert: []int{1, 2},
					delete: []int{1},
					x:      2,
					want:   true,
				},
				{
					insert: []int{1, 2},
					delete: []int{1},
					x:      1,
				},
			}

			for _, tc := range testCases {
				s := newSet()
				for _, in := range tc.insert {
					s.Insert(in)
				}
				for _, de := range tc.delete {
					s.Delete(de)
				}
				if got := s.Has(tc.x); got != tc.want {
					t.Errorf("After inserting %v and deleting %v, s.Has(%d) = %t, want %t", tc.insert, tc.delete, tc.x, got, tc.want)
				}
			}
		})
	}
}

func TestAllMinMax(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			testCases := []struct {
				insert    []int
				want      []int
				wantMin   int
				wantMinOK bool
				wantMax   int
				wantMaxOK bool
			}{
				{},
				{
					insert:    []int{1},
					want:      []int{1},
					wantMin:   1,
					wantMinOK: true,
					wantMax:   1,
					wantMaxOK: true,
				},
				{
					insert:    []int{1, 1},
					want:      []int{1},
					wantMin:   1,
					wantMinOK: true,
					wantMax:   1,
					wantMaxOK: true,
				},
				{
					insert:    []int{3, 2, 1},
					want:      []int{1, 2, 3},
					wantMin:   1,
					wantMinOK: true,
					wantMax:   3,
					wantMaxOK: true,
				},
				{
					insert:    []int{1, 2, 3, 2, 1},
					want:      []int{1, 2, 3},
					wantMin:   1,
					wantMinOK: true,
					wantMax:   3,
					wantMaxOK: true,
				},
			}

			for _, tc := range testCases {
				s := newSet()
				for _, in := range tc.insert {
					s.Insert(in)
				}
				var got []int
				for v := range s.All() {
					got = append(got, v)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("After calling Insert with each of %v, iterating All() produced %v, want %v.",
						tc.insert, got, tc.want)
				}
				if got, want := s.Len(), len(tc.want); got != want {
					t.Errorf("After calling Insert with each of %v, Len() = %v, want %v.",
						tc.insert, got, want)
				}
				if got, ok := s.Min(); got != tc.wantMin || ok != tc.wantMinOK {
					t.Errorf("After calling Insert with each of %v, Min() = %v, %t, want %v, %t.",
						tc.insert, got, ok, tc.wantMin, tc.wantMinOK)
				}
				if got, ok := s.Max(); got != tc.wantMax || ok != tc.wantMaxOK {
					t.Errorf("After calling Insert with each of %v, Max() = %v, %t, want %v, %t.",
						tc.insert, got, ok, tc.wantMax, tc.wantMaxOK)
				}
			}
		})
	}
}

type rangeTest struct {
	name string
	f    any
}

func testRange[T cmp.Ordered]() []rangeTest {
	var tests []rangeTest
	for name, newSet := range map[string]func() *Set[T]{
		"NewSet":     NewSet[T],
		"NewSetFunc": func() *Set[T] { return NewSetFunc(cmp.Compare[T]) },
	} {
		tests = append(tests, rangeTest{
			name: name,
			f: func(in []T) bool {
				s := newSet()
				for _, in := range in {
					s.Insert(in)
				}
				var got []T
				for v := range s.All() {
					got = append(got, v)
				}
				var gotBackward []T
				for v := range s.Backward() {
					gotBackward = append(gotBackward, v)
				}
				slices.Sort(in)
				want := slices.Compact(in)
				if !slices.Equal(got, want) {
					return false
				}
				slices.Reverse(want)
				return slices.Equal(gotBackward, want)
			},
		})
	}
	return tests
}

func TestRangeBlackBox(t *testing.T) {
	for name, tests := range map[string][]rangeTest{
		"int":     testRange[int](),
		"string":  testRange[string](),
		"float64": testRange[float64](),
		"byte":    testRange[byte](),
	} {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					if err := quick.Check(test.f, nil); err != nil {
						t.Error(err)
					}
				})
			}
		})
	}
}

func TestFindGreaterThanOrEqual(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			testCases := []struct {
				insert       []int
				x            int
				wantNotFound bool
				wantValue    int
			}{
				{
					x:            1,
					wantNotFound: true,
				},
				{
					x:            0,
					wantNotFound: true,
				},
				{
					insert:    []int{1},
					x:         1,
					wantValue: 1,
				},
				{
					insert:    []int{1},
					x:         0,
					wantValue: 1,
				},
				{
					insert:       []int{1},
					x:            2,
					wantNotFound: true,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         3,
					wantValue: 4,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         1,
					wantValue: 1,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         2,
					wantValue: 2,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         3,
					wantValue: 4,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         5,
					wantValue: 6,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         7,
					wantValue: 7,
				},
				{
					insert:    []int{2, 4, 6, 1, 7, 8},
					x:         8,
					wantValue: 8,
				},
				{
					insert:       []int{2, 4, 6, 1, 7, 8},
					x:            9,
					wantNotFound: true,
				},
			}

			for _, tc := range testCases {
				s := newSet()
				for _, in := range tc.insert {
					s.Insert(in)
				}
				si, ok := s.FindGreaterThanOrEqual(tc.x)
				if tc.wantNotFound {
					if ok || si != nil {
						t.Errorf("After inserting %v, s.FindGreaterThanOrEqual(%d) = %v, %t, want nil, false", tc.insert, tc.x, si, ok)
					}
				} else {
					if si == nil || si.Value() != tc.wantValue || !ok {
						t.Errorf("After inserting %v, s.FindGreaterThanOrEqual(%d) = %v, %t, want an item with value %d, true", tc.insert, tc.x, si, ok, tc.wantValue)
					}
				}
			}
		})
	}
}

func TestInsertNearby(t *testing.T) {
	for name, newSet := range map[string]func() *Set[int]{
		"NewSet":     NewSet[int],
		"NewSetFunc": func() *Set[int] { return NewSetFunc(cmp.Compare[int]) },
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("Increasing", func(t *testing.T) {
				s := newSet()
				if !s.Insert(1) {
					t.Errorf("newSet().Insert(1) = false, want true")
				}
				si, ok := s.First()
				if !ok {
					t.Fatalf("After Insert(1), First() returned false, want true")
				}
				for i := 2; i <= 10; i++ {
					si, ok = si.InsertNearby(i)
					if !ok {
						t.Fatalf("Inserting increasing sequence, InsertNearby(%d) returned false, want true", i)
					}
					if got := si.Value(); got != i {
						t.Errorf("InsertNearby(%d).Value() = %d, want %d", i, got, i)
					}
				}

				si, ok = si.InsertNearby(10)
				if ok {
					t.Fatalf("{item 10}.InsertNearby(10) returned true, want false")
				}
				if got := si.Value(); got != 10 {
					t.Errorf("InsertNearby(10).Value() = %d, want 10", got)
				}

				si, ok = si.InsertNearby(9)
				if ok {
					t.Fatalf("{item 10}.InsertNearby(9) returned true, want false")
				}
				if got := si.Value(); got != 9 {
					t.Errorf("InsertNearby(9).Value() = %d, want 9", got)
				}

				si, ok = si.InsertNearby(1)
				if ok {
					t.Fatalf("(item 9).InsertNearby(1) returned true, want false")
				}
				if got := si.Value(); got != 1 {
					t.Errorf("InsertNearby(1).Value() = %d, want 1", got)
				}

				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				if diff := gcmp.Diff(want, slices.Collect(s.All())); diff != "" {
					t.Errorf("After using InsertNearby with values from 1 to 10, All() diff (-want +got):\n%s", diff)
				}
			})
			t.Run("Decreasing", func(t *testing.T) {
				s := newSet()
				if !s.Insert(10) {
					t.Errorf("newSet().Insert(10) = false, want true")
				}
				si, ok := s.Last()
				if !ok {
					t.Fatalf("After Insert(10), Last() returned false, want true")
				}
				for i := 9; i >= 1; i-- {
					si, ok = si.InsertNearby(i)
					if !ok {
						t.Fatalf("Inserting decreasing sequence, InsertNearby(%d) returned false, want true", i)
					}
					if got := si.Value(); got != i {
						t.Errorf("InsertNearby(%d).Value() = %d, want %d", i, got, i)
					}
				}

				si, ok = si.InsertNearby(1)
				if ok {
					t.Fatalf("{item 1}.InsertNearby(1) returned true, want false")
				}
				if got := si.Value(); got != 1 {
					t.Errorf("InsertNearby(1).Value() = %d, want 1", got)
				}

				si, ok = si.InsertNearby(2)
				if ok {
					t.Fatalf("{item 1}.InsertNearby(2) returned true, want false")
				}
				if got := si.Value(); got != 2 {
					t.Errorf("InsertNearby(2).Value() = %d, want 2", got)
				}

				si, ok = si.InsertNearby(10)
				if ok {
					t.Fatalf("(item 2).InsertNearby(10) returned true, want false")
				}
				if got := si.Value(); got != 10 {
					t.Errorf("InsertNearby(10).Value() = %d, want 10", got)
				}

				want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				if diff := gcmp.Diff(want, slices.Collect(s.All())); diff != "" {
					t.Errorf("After using InsertNearby with values from 10 to 1, All() diff (-want +got):\n%s", diff)
				}
			})
			t.Run("Between", func(t *testing.T) {
				t.Run("AfterFirst", func(t *testing.T) {
					s := newSet()
					s.Insert(1)
					s.Insert(3)
					first, ok := s.First()
					if !ok {
						t.Fatalf("After inserting 1 and 3, there is still no First() item.")
					}
					if got, ok := first.InsertNearby(2); !ok || got.Value() != 2 {
						t.Errorf("first.InsertNearby(2) = %s, %t, want {value 2}, true", got, ok)
					}
				})
				t.Run("BeforeLast", func(t *testing.T) {
					s := newSet()
					s.Insert(1)
					s.Insert(3)
					last, ok := s.Last()
					if !ok {
						t.Fatalf("After inserting 1 and 3, there is still no Last() item.")
					}
					if got, ok := last.InsertNearby(2); !ok || got.Value() != 2 {
						t.Errorf("last.InsertNearby(2) = %s, %t, want {value 2}, true", got, ok)
					}
				})
			})
		})
	}
}

func (s *SetItem[T]) String() string {
	if s == nil {
		return "nil"
	}
	if s.level == 0 {
		return "bottom"
	}
	return fmt.Sprintf("{value: %v, level: %d, left: %v, right: %v}", s.value, s.level, s.l, s.r)
}
