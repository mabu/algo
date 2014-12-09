// Package maxflow provides algorithms that solve maximum flow problem.
package maxflow

import (
	"github.com/mabu/algo/graph"
	"github.com/mabu/algo/graph/bfs"
)

// Graph specifies the graph interface required by maxflow algorithms.
type Graph interface {
	graph.Graph
	// AddFlow increases the flow from v to u by the amount, and
	// decreases the flow from u to v by the same amount.
	AddFlow(v, u, amount int)
	// ResidualCapacity tells the remaining capacity from v to u.
	ResidualCapacity(v, u int) int
}

// EdmondsKarp runs the Edmonds–Karp algorithm and returns the
// maximum flow from source to sink. If source and sink is the same,
// panics.
func EdmondsKarp(g Graph, source, sink int) int {
	flow := 0
	for {
		path := bfs.Path(residualGraph{g}, source, sink)
		if len(path) == 0 {
			break
		}
		if len(path) == 1 {
			panic("source == sink")
		}
		f := g.ResidualCapacity(path[0], path[1])
		for i := range path[2:] {
			if r := g.ResidualCapacity(path[i+1], path[i+2]); r < f {
				f = r
			}
		}
		if f <= 0 {
			panic("invalid graph")
		}
		for i := range path[1:] {
			g.AddFlow(path[i], path[i+1], f)
		}
		flow += f
	}
	return flow
}

type residualGraph struct {
	Graph
}

func (g residualGraph) Adjacent(v int) []int {
	var res []int
	for _, u := range g.Graph.Adjacent(v) {
		if g.ResidualCapacity(v, u) > 0 {
			res = append(res, u)
		}
	}
	return res
}
