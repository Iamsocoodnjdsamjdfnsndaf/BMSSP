package dijkstra

import (
	"math"
	"playground/common"
	"testing"
)

func TestDijkstra_SimpleGraph(t *testing.T) {
	g := createLinearGraph(5)
	algo := NewDijkstraAlgorithm(g, []int{0}, nil)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	expected := map[int]float64{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstra_MultiSource(t *testing.T) {
	g := createLinearGraph(5)
	sources := []int{0, 4}
	algo := NewDijkstraAlgorithm(g, sources, nil)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	expected := map[int]float64{0: 0, 1: 1, 2: 2, 3: 1, 4: 0}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstra_Bounded(t *testing.T) {
	g := createLinearGraph(10)
	boundary := 5.0
	algo := NewDijkstraAlgorithm(g, []int{0}, &boundary)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	for i := 0; i < 5; i++ {
		if dist[i] != float64(i) {
			t.Errorf("Vertex %d: expected dist=%f, got %f", i, float64(i), dist[i])
		}
	}
	for i := 5; i < 10; i++ {
		if !math.IsInf(dist[i], 1) {
			t.Errorf("Vertex %d beyond boundary should have infinite distance, got %f", i, dist[i])
		}
	}
}

func TestDijkstra_DisconnectedGraph(t *testing.T) {
	g := &common.Graph{N: 6, Adj: make(map[int][]common.Edge)}
	edges := []common.Edge{
		{U: 0, V: 1, Weight: 1.0}, {U: 1, V: 2, Weight: 1.0},
		{U: 3, V: 4, Weight: 1.0}, {U: 4, V: 5, Weight: 1.0},
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}

	algo := NewDijkstraAlgorithm(g, []int{0}, nil)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	if math.IsInf(dist[1], 1) || math.IsInf(dist[2], 1) {
		t.Error("Connected vertices have infinite distance")
	}
	if !math.IsInf(dist[3], 1) || !math.IsInf(dist[4], 1) || !math.IsInf(dist[5], 1) {
		t.Error("Disconnected vertices have finite distance")
	}
}

func TestDijkstra_CycleGraph(t *testing.T) {
	g := createCycleGraph()
	algo := NewDijkstraAlgorithm(g, []int{0}, nil)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestDijkstra_WeightedCompleteGraph(t *testing.T) {
	g := createWeightedCompleteGraph()
	algo := NewDijkstraAlgorithm(g, []int{0}, nil)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

// --- Benchmarks ---

func BenchmarkDijkstra_Small(b *testing.B) {
	g := createLinearGraph(100)
	sources := []int{0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		algo := NewDijkstraAlgorithm(g, sources, nil)
		_, _ = algo.Solve()
	}
}

func BenchmarkDijkstra_Medium(b *testing.B) {
	g := createLinearGraph(1000)
	sources := []int{0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		algo := NewDijkstraAlgorithm(g, sources, nil)
		_, _ = algo.Solve()
	}
}

func BenchmarkDijkstra_Large(b *testing.B) {
	g := createLinearGraph(10000)
	sources := []int{0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		algo := NewDijkstraAlgorithm(g, sources, nil)
		_, _ = algo.Solve()
	}
}

func BenchmarkDijkstra_MultiSource(b *testing.B) {
	g := createLinearGraph(1000)
	sources := []int{0, 250, 500, 750, 999}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		algo := NewDijkstraAlgorithm(g, sources, nil)
		_, _ = algo.Solve()
	}
}

func BenchmarkDijkstra_Bounded(b *testing.B) {
	g := createLinearGraph(10000)
	sources := []int{0}
	boundary := 100.0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		algo := NewDijkstraAlgorithm(g, sources, &boundary)
		_, _ = algo.Solve()
	}
}

// --- Helper Functions ---

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

func createCycleGraph() *common.Graph {
	g := &common.Graph{N: 4, Adj: make(map[int][]common.Edge)}
	edges := []common.Edge{
		{U: 0, V: 1, Weight: 1.0}, {U: 1, V: 2, Weight: 2.0},
		{U: 2, V: 3, Weight: 1.0}, {U: 3, V: 0, Weight: 5.0},
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges
	return g
}

func createWeightedCompleteGraph() *common.Graph {
	g := &common.Graph{N: 4, Adj: make(map[int][]common.Edge)}
	weights := [][]float64{
		{0, 1, 3, 7}, {1, 0, 2, 5},
		{3, 2, 0, 1}, {7, 5, 1, 0},
	}
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			w := weights[i][j]
			g.Adj[i] = append(g.Adj[i], common.Edge{U: i, V: j, Weight: w})
			g.Adj[j] = append(g.Adj[j], common.Edge{U: j, V: i, Weight: w})
		}
	}
	return g
}
