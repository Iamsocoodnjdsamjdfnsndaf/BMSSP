package dijkstra

import (
	"math"
	"playground/common"
	"testing"
)

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

func TestDijkstra_SimpleGraph(t *testing.T) {
	g := createLinearGraph(5)
	dist := Dijkstra(g, 0)

	expected := map[int]float64{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstraMultiSource(t *testing.T) {
	g := createLinearGraph(5)
	dist := DijkstraMultiSource(g, []int{0, 4})

	// Middle vertex should have distance 2 from nearest source
	if dist[2] != 2.0 {
		t.Errorf("Middle vertex distance incorrect: expected 2, got %f", dist[2])
	}

	// Check all distances
	expected := map[int]float64{0: 0, 1: 1, 2: 2, 3: 1, 4: 0}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstra_DisconnectedGraph(t *testing.T) {
	g := &common.Graph{
		N:   6,
		Adj: make(map[int][]common.Edge),
	}

	// Two disconnected components
	edges := []common.Edge{
		{U: 0, V: 1, Weight: 1.0},
		{U: 1, V: 2, Weight: 1.0},
		{U: 3, V: 4, Weight: 1.0},
		{U: 4, V: 5, Weight: 1.0},
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}

	dist := Dijkstra(g, 0)

	// Vertices in same component should have finite distance
	if math.IsInf(dist[1], 1) || math.IsInf(dist[2], 1) {
		t.Error("Connected vertices have infinite distance")
	}

	// Vertices in other component should have infinite distance
	if !math.IsInf(dist[3], 1) || !math.IsInf(dist[4], 1) || !math.IsInf(dist[5], 1) {
		t.Error("Disconnected vertices have finite distance")
	}
}

func TestDijkstra_CycleGraph(t *testing.T) {
	g := &common.Graph{
		N:   4,
		Adj: make(map[int][]common.Edge),
	}

	// Create cycle with different weights
	edges := []common.Edge{
		{U: 0, V: 1, Weight: 1.0},
		{U: 1, V: 2, Weight: 2.0},
		{U: 2, V: 3, Weight: 1.0},
		{U: 3, V: 0, Weight: 5.0}, // Completes cycle
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges

	dist := Dijkstra(g, 0)

	// Check shortest paths
	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstraBounded(t *testing.T) {
	g := createLinearGraph(10)
	B := 5.0 // Boundary at distance 5

	dist := DijkstraBounded(g, []int{0}, B)

	// Vertices within boundary should have correct distances
	for i := 0; i < 5; i++ {
		if dist[i] != float64(i) {
			t.Errorf("Vertex %d: expected dist=%f, got %f", i, float64(i), dist[i])
		}
	}

	// Vertices beyond boundary should have infinite distance
	for i := 5; i < 10; i++ {
		if !math.IsInf(dist[i], 1) {
			t.Errorf("Vertex %d beyond boundary should have infinite distance, got %f", i, dist[i])
		}
	}
}

func TestDijkstra_WeightedCompleteGraph(t *testing.T) {
	// Complete graph K4
	g := &common.Graph{
		N:   4,
		Adj: make(map[int][]common.Edge),
	}

	// All pairs connected with varying weights
	weights := [][]float64{
		{0, 1, 3, 7},
		{1, 0, 2, 5},
		{3, 2, 0, 1},
		{7, 5, 1, 0},
	}

	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			w := weights[i][j]
			g.Adj[i] = append(g.Adj[i], common.Edge{U: i, V: j, Weight: w})
			g.Adj[j] = append(g.Adj[j], common.Edge{U: j, V: i, Weight: w})
		}
	}

	dist := Dijkstra(g, 0)

	// Verify shortest paths
	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

// Benchmarks for Dijkstra alone
func BenchmarkDijkstra_Small(b *testing.B) {
	g := createLinearGraph(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dijkstra(g, 0)
	}
}

func BenchmarkDijkstra_Medium(b *testing.B) {
	g := createLinearGraph(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dijkstra(g, 0)
	}
}

func BenchmarkDijkstra_Large(b *testing.B) {
	g := createLinearGraph(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dijkstra(g, 0)
	}
}

func BenchmarkDijkstraMultiSource(b *testing.B) {
	g := createLinearGraph(1000)
	sources := []int{0, 250, 500, 750, 999}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DijkstraMultiSource(g, sources)
	}
}

func BenchmarkDijkstraBounded(b *testing.B) {
	g := createLinearGraph(10000)
	sources := []int{0}
	B := 100.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DijkstraBounded(g, sources, B)
	}
}
