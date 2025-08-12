package benchmarks

import (
	"fmt"
	"math/rand"
	"playground/bmssp"
	"playground/common"
	"playground/dijkstra"
	"testing"
	"time"
)

// --- Benchmark Comparisons ---

func BenchmarkComparison_LinearGraph_Small(b *testing.B) {
	g := createLinearGraph(100)
	sources := []int{0}
	B := 1000
	l := 3

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_LinearGraph_Medium(b *testing.B) {
	g := createLinearGraph(1000)
	sources := []int{0}
	B := 10000
	l := 4

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_LinearGraph_Large(b *testing.B) {
	g := createLinearGraph(10000)
	sources := []int{0}
	B := 100000
	l := 5

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_GridGraph_Small(b *testing.B) {
	g := createGridGraph(10, 10) // 100 vertices
	sources := []int{0}
	B := 1000
	l := 3

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_GridGraph_Medium(b *testing.B) {
	g := createGridGraph(50, 50) // 2500 vertices
	sources := []int{0}
	B := 10000
	l := 4

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_RandomGraph_Sparse(b *testing.B) {
	n := 1000
	m := 2000 // Sparse: ~2 edges per vertex
	g := createRandomGraph(n, m, time.Now().Unix())
	sources := []int{0}
	B := 10000
	l := 4

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_RandomGraph_Dense(b *testing.B) {
	n := 500
	m := 10000 // Dense: ~20 edges per vertex
	g := createRandomGraph(n, m, time.Now().Unix())
	sources := []int{0}
	B := 10000
	l := 4

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_MultiSource(b *testing.B) {
	g := createLinearGraph(1000)
	sources := []int{0, 250, 500, 750, 999}
	B := 10000
	l := 4

	b.Run("BMSSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra_MultiSource", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkComparison_Bounded(b *testing.B) {
	g := createGridGraph(100, 100) // 10000 vertices
	sources := []int{0}
	boundary := 50.0
	l := 3

	b.Run("BMSSP_Bounded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, boundary, sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra_Bounded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, &boundary)
			_, _ = algo.Solve()
		}
	})
}

func BenchmarkMemory_Comparison(b *testing.B) {
	g := createGridGraph(50, 50) // 2500 vertices
	sources := []int{0}
	B := 10000
	l := 4

	b.Run("BMSSP_Memory", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			algo := bmssp.NewBMSSPAlgorithm(g, l, float64(B), sources)
			_, _ = algo.Solve()
		}
	})

	b.Run("Dijkstra_Memory", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			algo := dijkstra.NewDijkstraAlgorithm(g, sources, nil)
			_, _ = algo.Solve()
		}
	})
}

// --- Graph Generation Utilities ---

func createLinearGraph(n int) *common.Graph {
	g := &common.Graph{
		N:   n,
		Adj: make(map[int][]common.Edge),
	}
	edges := make([]common.Edge, 0, n-1)
	for i := 0; i < n-1; i++ {
		edges = append(edges, common.Edge{U: i, V: i + 1, Weight: 1.0})
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges
	return g
}

func createGridGraph(rows, cols int) *common.Graph {
	n := rows * cols
	g := &common.Graph{
		N:   n,
		Adj: make(map[int][]common.Edge),
	}
	edges := make([]common.Edge, 0)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := i*cols + j
			if j < cols-1 {
				u := i*cols + (j + 1)
				edges = append(edges, common.Edge{U: v, V: u, Weight: 1.0})
			}
			if i < rows-1 {
				u := (i+1)*cols + j
				edges = append(edges, common.Edge{U: v, V: u, Weight: 1.0})
			}
		}
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges
	return g
}

func createRandomGraph(n, m int, seed int64) *common.Graph {
	r := rand.New(rand.NewSource(seed))
	g := &common.Graph{
		N:   n,
		Adj: make(map[int][]common.Edge),
	}
	edges := make([]common.Edge, 0)
	edgeSet := make(map[string]bool)

	for i := 1; i < n; i++ {
		parent := r.Intn(i)
		weight := r.Float64()*10 + 1
		edges = append(edges, common.Edge{U: parent, V: i, Weight: weight})
		edgeSet[fmt.Sprintf("%d-%d", parent, i)] = true
		edgeSet[fmt.Sprintf("%d-%d", i, parent)] = true
	}

	for len(edges) < m && len(edges) < n*(n-1)/2 {
		u := r.Intn(n)
		v := r.Intn(n)
		if u != v {
			key1 := fmt.Sprintf("%d-%d", u, v)
			key2 := fmt.Sprintf("%d-%d", v, u)
			if !edgeSet[key1] {
				weight := r.Float64()*10 + 1
				edges = append(edges, common.Edge{U: u, V: v, Weight: weight})
				edgeSet[key1] = true
				edgeSet[key2] = true
			}
		}
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges
	return g
}
