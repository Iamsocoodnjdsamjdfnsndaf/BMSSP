package main

import (
	"fmt"
	"math"
	"playground/bmssp"
	"playground/common"
)

func main() {
	g := &common.Graph{
		N:   5,
		Adj: make(map[int][]common.Edge),
	}

	addUndirected := func(u, v int, w float64) {
		e1 := common.Edge{U: u, V: v, Weight: w}
		e2 := common.Edge{U: v, V: u, Weight: w}
		g.Adj[u] = append(g.Adj[u], e1)
		g.Adj[v] = append(g.Adj[v], e2)
		g.Edges = append(g.Edges, e1, e2)
	}

	addUndirected(0, 1, 1.0)
	addUndirected(1, 2, 2.0)
	addUndirected(2, 3, 1.0)
	addUndirected(3, 4, 3.0)

	S := []int{0} // sources
	B := 1000.0   // boundary must be float64
	l := 3        // recursion depth

	algo := bmssp.NewBMSSPAlgorithm(g, l, B, S)
	dist, err := algo.Solve()
	if err != nil {
		panic(err)
	}

	for v := 0; v < g.N; v++ {
		if math.IsInf(dist[v], 1) {
			fmt.Printf("dist[%d] = +Inf\n", v)
		} else {
			fmt.Printf("dist[%d] = %.6f\n", v, dist[v])
		}
	}
}
