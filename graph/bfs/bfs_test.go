package bfs

import (
	"reflect"
	"testing"

	"github.com/mabu/algo/graph"
)

type unsized [][]int

func (g unsized) Adjacent(v int) []int {
	return g[v]
}

func TestPath(t *testing.T) {
	for _, tc := range []struct {
		g        [][]int
		from, to int
		want     []int
	}{
		{want: []int{0}},
		{
			g: [][]int{
				{1},
				{},
			},
			from: 0,
			to:   1,
			want: []int{0, 1},
		},
		{
			g: [][]int{
				{1},
				{},
			},
			from: 0,
			to:   0,
			want: []int{0},
		},
		{
			g: [][]int{
				{1},
				{},
			},
			from: 1,
			to:   1,
			want: []int{1},
		},
		{
			g: [][]int{
				{1},
				{2},
				{3},
				{0},
			},
			from: 2,
			to:   1,
			want: []int{2, 3, 0, 1},
		},
		{
			g: [][]int{
				{1, 2},
				{1, 3},
				{4},
				{1, 2, 4},
				{1, 0, 4, 2, 3},
			},
			from: 0,
			to:   4,
			want: []int{0, 2, 4},
		},
	} {
		if got := Path(unsized(tc.g), tc.from, tc.to); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Unsized Path(%v, %d, %d) = %v, want %v", tc.g, tc.from, tc.to, got, tc.want)
		}
		if got := Path(graph.AdjList(tc.g), tc.from, tc.to); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Sized Path(%v, %d, %d) = %v, want %v", tc.g, tc.from, tc.to, got, tc.want)
		}
	}
}
