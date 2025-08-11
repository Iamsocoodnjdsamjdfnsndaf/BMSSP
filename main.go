package main

import (
	"math"
	"playground/bmssp"
)

func main() {
	// Example usage
	g := &bmssp.Graph{
		N:   5,
		Adj: make(map[int][]bmssp.Edge),
	}

	// Build adjacency list
	edges := []bmssp.Edge{
		{U: 0, V: 1, Weight: 1.0},
		{U: 1, V: 2, Weight: 2.0},
		{U: 2, V: 3, Weight: 1.0},
		{U: 3, V: 4, Weight: 3.0},
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], bmssp.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges

	S := []int{0} // Source vertices
	B := 100      // Boundary
	l := 3        // Recursion depth

	dist := make(map[int]float64)
	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	for _, s := range S {
		dist[s] = 0
	}

	_, _ = bmssp.BMSSP(l, B, S, g, dist)
}
