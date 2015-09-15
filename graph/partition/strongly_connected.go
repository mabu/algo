// Package partition contains graph partitioning algorithms.
package partition

import "github.com/mabu/algo/graph"

// StronglyConnected returns strongly connected components of the graph.
func StronglyConnected(g graph.Sized) [][]int {
	// Tarjan's algorithm.
	type vertex struct {
		onStack    bool
		low, index int
	}
	vs := make([]vertex, g.Size())
	nextIndex := 1
	var stack []int
	var components [][]int
	var dfs func(int)
	dfs = func(v int) {
		vs[v].index = nextIndex
		vs[v].low = nextIndex
		nextIndex++
		stack = append(stack, v)
		vs[v].onStack = true

		for _, w := range g.Adjacent(v) {
			if vs[w].index == 0 {
				dfs(w)
				if l := vs[w].low; vs[v].low > l {
					vs[v].low = l
				}
			} else if vs[w].onStack {
				if l := vs[w].low; vs[v].low > l {
					vs[v].low = l
				}
			}
		}

		if vs[v].low == vs[v].index {
			i := len(stack) - 1
			for stack[i] != v {
				i--
			}
			comp := make([]int, len(stack)-i)
			copy(comp, stack[i:])
			components = append(components, comp)
			stack = stack[:i]
			for _, u := range comp {
				vs[u].onStack = false
			}
		}
	}

	for i := range vs {
		if vs[i].index == 0 {
			dfs(i)
		}
	}
	return components
}
