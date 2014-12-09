// Package bfs implements breadth-first search algorithm.
package bfs

import "github.com/mabu/algo/graph"

// Path returns one of the shortest paths, including both endpoints.
// If the path does not exist, returns nil.
// May be more efficient if g implements graph.Sized.
func Path(g graph.Graph, from, to int) []int {
	if from == to {
		return []int{from}
	}
	var parent parent
	if s, ok := g.(graph.Sized); ok {
		parent = newParentSlice(s.Size())
	} else {
		parent = make(parentMap)
	}
	parent.set(from, from)
	queue := []int{from}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, u := range g.Adjacent(v) {
			if u == to {
				var reversed []int
				for w := v; w != from; w = parent.get(w) {
					reversed = append(reversed, w)
				}
				rl := len(reversed)
				path := make([]int, rl+2)
				path[0] = from
				path[rl+1] = to
				for i, r := range reversed {
					path[rl-i] = r
				}
				return path
			}
			if parent.has(u) {
				continue
			}
			parent.set(u, v)
			queue = append(queue, u)
		}
	}
	return nil
}

type parent interface {
	has(int) bool
	get(int) int
	set(c, p int)
}

type parentMap map[int]int

func (m parentMap) has(v int) bool {
	_, ok := m[v]
	return ok
}

func (m parentMap) get(v int) int {
	return m[v]
}

func (m parentMap) set(c, p int) {
	m[c] = p
}

type parentSlice struct {
	isSet []bool
	p     []int
}

func newParentSlice(size int) parentSlice {
	return parentSlice{make([]bool, size), make([]int, size)}
}

func (s parentSlice) has(v int) bool {
	return s.isSet[v]
}

func (s parentSlice) get(v int) int {
	return s.p[v]
}

func (s parentSlice) set(c, p int) {
	s.p[c] = p
	s.isSet[c] = true
}
