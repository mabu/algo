package partition

import (
	"reflect"
	"sort"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mabu/algo/graph"
)

func TestStronglyConnected(t *testing.T) {
	for _, tc := range []struct {
		graph graph.AdjList
		want  [][]int
	}{
		{},
		{
			graph.AdjList{
				0: {},
			},
			[][]int{{0}},
		},
		{
			graph.AdjList{
				0: {0},
			},
			[][]int{{0}},
		},
		{
			graph.AdjList{
				0: {1},
				1: {},
			},
			[][]int{{0}, {1}},
		},
		{
			graph.AdjList{
				0: {1},
				1: {0},
			},
			[][]int{{0, 1}},
		},
		{
			graph.AdjList{
				0: {1, 2},
				1: {2},
				2: {},
			},
			[][]int{{0}, {1}, {2}},
		},
		{
			graph.AdjList{
				0: {1, 2},
				1: {2},
				2: {0},
			},
			[][]int{{0, 1, 2}},
		},
		{
			graph.AdjList{
				0: {1},
				1: {2},
				2: {0},
			},
			[][]int{{0, 1, 2}},
		},
		{
			graph.AdjList{
				0: {2},
				1: {2},
				2: {0},
			},
			[][]int{{0, 2}, {1}},
		},
		{
			graph.AdjList{
				0: {1},
				1: {2},
				2: {3},
				3: {1, 4},
				4: {},
			},
			[][]int{{0}, {1, 2, 3}, {4}},
		},
		{
			graph.AdjList{
				0: {1},
				1: {2},
				2: {3},
				3: {1, 4},
				4: {0},
			},
			[][]int{{0, 1, 2, 3, 4}},
		},
	} {
		got := StronglyConnected(tc.graph)
		sortComponents(got)
		sortComponents(tc.want)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Graph %v: diff -want +got:\n%s", tc.graph, pretty.Compare(tc.want, got))
		}
	}
}

func sortComponents(cs [][]int) {
	for _, c := range cs {
		sort.Ints(c)
	}
	sort.Sort(lexicographically(cs))
}

type lexicographically [][]int

func (l lexicographically) Less(i, j int) bool {
	for k := range l[i] {
		if k >= len(l[j]) {
			return false
		}
		if l[i][k] < l[j][k] {
			return true
		}
		if l[i][k] > l[j][k] {
			return false
		}
	}
	return false
}

func (l lexicographically) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l lexicographically) Len() int      { return len(l) }
