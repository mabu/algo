package maxflow

import (
	"reflect"
	"testing"
)

type edge struct {
	from, to int
}

type testGraph struct {
	edges          [][]int
	capacity, flow map[edge]int
}

func (g testGraph) Adjacent(v int) []int {
	return g.edges[v]
}

func (g testGraph) AddFlow(v, u, amount int) {
	g.flow[edge{v, u}] += amount
	g.flow[edge{u, v}] -= amount
}

func (g testGraph) ResidualCapacity(v, u int) int {
	e := edge{v, u}
	return g.capacity[e] - g.flow[e]
}

func TestEdmondsKarp(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		edges        [][]struct{ to, capacity int }
		source, sink int
		wantFlow     map[edge]int
		want         int
	}{
		{
			desc: "trivial",
			edges: [][]struct{ to, capacity int }{
				{{1, 4}},
				{{0, 0}},
			},
			source: 0,
			sink:   1,
			wantFlow: map[edge]int{
				edge{0, 1}: 4,
				edge{1, 0}: -4,
			},
			want: 4,
		},
		{
			desc: "http://upload.wikimedia.org/wikipedia/commons/thumb/9/94/Max_flow.svg/330px-Max_flow.svg.png",
			edges: [][]struct{ to, capacity int }{
				{{1, 3}, {4, 3}},         // s
				{{0, 0}, {4, 2}, {2, 3}}, // o
				{{1, 0}, {5, 4}, {3, 2}}, // q
				{{2, 0}, {5, 0}},         // t
				{{0, 0}, {1, 0}, {5, 2}}, // p
				{{4, 0}, {2, 0}, {3, 3}}, // r
			},
			source: 0,
			sink:   3,
			wantFlow: map[edge]int{
				edge{0, 1}: 3,
				edge{1, 0}: -3,
				edge{0, 4}: 2,
				edge{4, 0}: -2,
				edge{1, 2}: 3,
				edge{2, 1}: -3,
				edge{2, 3}: 2,
				edge{3, 2}: -2,
				edge{2, 5}: 1,
				edge{5, 2}: -1,
				edge{4, 5}: 2,
				edge{5, 4}: -2,
				edge{5, 3}: 3,
				edge{3, 5}: -3,
			},
			want: 5,
		},
		{
			desc: "http://upload.wikimedia.org/wikipedia/commons/thumb/3/3e/Edmonds-Karp_flow_example_0.svg/300px-Edmonds-Karp_flow_example_0.svg.png",
			edges: [][]struct{ to, capacity int }{
				{{1, 3}, {2, 0}, {3, 3}},
				{{0, 0}, {2, 4}, {4, 0}},
				{{0, 3}, {1, 0}, {3, 1}, {4, 2}},
				{{0, 0}, {2, 0}, {4, 2}, {5, 6}},
				{{1, 1}, {2, 0}, {3, 0}, {6, 1}},
				{{3, 0}, {6, 9}},
				{{4, 0}, {5, 0}},
			},
			source: 0,
			sink:   6,
			wantFlow: map[edge]int{
				edge{0, 3}: 3,
				edge{3, 0}: -3,
				edge{3, 5}: 4,
				edge{5, 3}: -4,
				edge{0, 1}: 2,
				edge{1, 0}: -2,
				edge{2, 3}: 1,
				edge{3, 2}: -1,
				edge{3, 4}: 0,
				edge{4, 3}: 0,
				edge{5, 6}: 4,
				edge{6, 5}: -4,
				edge{1, 2}: 2,
				edge{2, 1}: -2,
				edge{2, 4}: 1,
				edge{4, 2}: -1,
				edge{4, 6}: 1,
				edge{6, 4}: -1,
			},
			want: 5,
		},
	} {
		g := testGraph{
			edges:    make([][]int, len(tc.edges)),
			capacity: make(map[edge]int),
			flow:     make(map[edge]int),
		}
		for v := range tc.edges {
			for _, e := range tc.edges[v] {
				g.edges[v] = append(g.edges[v], e.to)
				g.capacity[edge{v, e.to}] = e.capacity
			}
		}
		t.Logf("Running %s: EdmondsKarp(%v, %d, %d)", tc.desc, g, tc.source, tc.sink)
		if got := EdmondsKarp(g, tc.source, tc.sink); got != tc.want {
			t.Errorf("got total flow %v, want %v", got, tc.want)
		}
		if !reflect.DeepEqual(g.flow, tc.wantFlow) {
			t.Errorf("got flow %v, want %v", g.flow, tc.wantFlow)
		}
	}
}
