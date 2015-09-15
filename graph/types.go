// Package graph specifies interfaces for graphs with integer-labeled nodes.
package graph

// Graph is the most general interface for a graph.
type Graph interface {
	// Adjacent reports what nodes are adjacent to the node v.
	Adjacent(v int) []int
}

// Sized is a graph with a known size and nodes enumerated from 0 to Size() - 1.
// Passing a type that implements Sized may improve algorithm performance even
// though they also work with a general Graph.
type Sized interface {
	Graph
	// Size is the number of nodes in the graph.
	Size() int
}

// AdjList is a graph represented as an adjacency list. Implements Sized.
type AdjList [][]int

func (l AdjList) Adjacent(v int) []int { return l[v] }
func (l AdjList) Size() int            { return len(l) }
