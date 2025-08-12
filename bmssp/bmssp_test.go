package bmssp

import (
	"math"
	"playground/common"
	"testing"
)

func TestDataStructureD_BasicOperations(t *testing.T) {
	var d DataStructureD
	d.Initialize(2, 100.0) // M=2, B=100.0

	if !d.IsEmpty() {
		t.Error("Expected empty data structure after initialization")
	}

	d.Insert(1, 5.0)
	d.Insert(2, 3.0)
	d.Insert(3, 7.0)
	d.Insert(4, 1.0)

	if d.IsEmpty() {
		t.Error("Expected non-empty after insertions")
	}

	B1, S1 := d.Pull()
	if B1 != 5.0 {
		t.Errorf("Pull 1: Expected boundary 5.0, got %f", B1)
	}
	if len(S1) != 2 || S1[0] != 4 || S1[1] != 2 {
		t.Errorf("Pull 1: Expected S={4, 2}, got %v", S1)
	}

	B2, S2 := d.Pull()
	if B2 != 100.0 {
		t.Errorf("Pull 2: Expected boundary 100.0, got %f", B2)
	}
	if len(S2) != 2 || S2[0] != 1 || S2[1] != 3 {
		t.Errorf("Pull 2: Expected S={1, 3}, got %v", S2)
	}

	if !d.IsEmpty() {
		t.Error("Expected empty data structure after two pulls")
	}
}

func TestDataStructureD_BatchPrepend(t *testing.T) {
	var d DataStructureD
	d.Initialize(10, 100.0)

	entries := []common.DistEntry{
		{Vertex: 1, Dist: 5.0},
		{Vertex: 2, Dist: 3.0},
		{Vertex: 3, Dist: 7.0},
		{Vertex: 2, Dist: 2.0}, // Duplicate with smaller value
	}

	d.BatchPrepend(entries)

	if d.pq.Len() != 3 {
		t.Errorf("Expected 3 elements after batch prepend, got %d", d.pq.Len())
	}

	_, S := d.Pull()
	expected := []int{2, 1, 3}
	if len(S) != len(expected) {
		t.Fatalf("Expected %d pulled elements, got %d", len(expected), len(S))
	}
	for i := range expected {
		if S[i] != expected[i] {
			t.Errorf("Unexpected pull order: expected %v, got %v", expected, S)
		}
	}
}

func TestDataStructureD_TieDrain_NoSplit(t *testing.T) {
	var d DataStructureD
	d.Initialize(1, 1000) // M=1 on purpose

	// three equal keys, then a larger one
	d.Insert(10, 0)
	d.Insert(11, 0)
	d.Insert(12, 0)
	d.Insert(20, 5)

	Bi, S := d.Pull()
	if len(S) != 3 {
		t.Fatalf("expected to drain all ties, got %d elems: %v", len(S), S)
	}
	if Bi != 5 {
		t.Fatalf("expected Bi=5 (next strictly larger), got %v", Bi)
	}
}

func TestBMSSP_FractionalWeights(t *testing.T) {
	g := &common.Graph{N: 6, Adj: make(map[int][]common.Edge)}
	add := func(u, v int, w float64) {
		g.Adj[u] = append(g.Adj[u], common.Edge{U: u, V: v, Weight: w})
		g.Adj[v] = append(g.Adj[v], common.Edge{U: v, V: u, Weight: w})
	}
	for i := 0; i < 5; i++ {
		add(i, i+1, 0.6)
	}

	algo := NewBMSSPAlgorithm(g, 3, 1000, []int{0})
	dist, _ := algo.Solve()

	for i := 0; i < 6; i++ {
		want := float64(i) * 0.6
		if math.Abs(dist[i]-want) > 1e-9 {
			t.Fatalf("v%d: want %.1f, got %v", i, want, dist[i])
		}
	}
}

func TestBMSSP_NoPivots_Successful(t *testing.T) {
	// line of 6 with weight 1; B so tight no relax happens
	g := &common.Graph{N: 6, Adj: make(map[int][]common.Edge)}
	for i := 0; i < 5; i++ {
		g.Adj[i] = append(g.Adj[i], common.Edge{U: i, V: i + 1, Weight: 1})
		g.Adj[i+1] = append(g.Adj[i+1], common.Edge{U: i + 1, V: i, Weight: 1})
	}
	algo := NewBMSSPAlgorithm(g, 3, 0.5, []int{0})
	dist, _ := algo.Solve()

	if dist[0] != 0 || !math.IsInf(dist[1], 1) {
		t.Fatalf("expected only source reached under tight B; got dist=%v", dist)
	}
}

// --- Tests for BMSSPAlgorithm ---

func TestBMSSP_SingleSource(t *testing.T) {
	g := createLinearGraph(5)
	algo := NewBMSSPAlgorithm(g, 2, 100.0, []int{0})
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	expectedDist := map[int]float64{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	for v, expected := range expectedDist {
		if math.Abs(dist[v]-expected) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, expected, dist[v])
		}
	}
}

func TestBMSSP_MultiSource(t *testing.T) {
	g := createLinearGraph(5)
	S := []int{0, 4} // Sources at both ends
	algo := NewBMSSPAlgorithm(g, 2, 100.0, S)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	// Middle vertex should have distance 2 from nearest source
	if dist[2] != 2.0 {
		t.Errorf("Middle vertex distance incorrect: expected 2, got %f", dist[2])
	}
}

func TestBMSSP_DisconnectedGraph(t *testing.T) {
	g := &common.Graph{
		N:   6,
		Adj: make(map[int][]common.Edge),
	}
	edges := []common.Edge{
		{0, 1, 1.0}, {1, 2, 1.0}, // Component 1
		{3, 4, 1.0}, {4, 5, 1.0}, // Component 2
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges

	algo := NewBMSSPAlgorithm(g, 2, 100.0, []int{0})
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

func TestBMSSP_CycleGraph(t *testing.T) {
	g := createCycleGraph()
	algo := NewBMSSPAlgorithm(g, 2, 100.0, []int{0})
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

func TestBMSSP_DifferentRecursionDepths(t *testing.T) {
	g := createLinearGraph(10)
	for l := 0; l <= 4; l++ {
		algo := NewBMSSPAlgorithm(g, l, 100.0, []int{0})
		dist, err := algo.Solve()
		if err != nil {
			t.Fatalf("l=%d: Solve() returned an error: %v", l, err)
		}

		for i := 0; i < g.N; i++ {
			if math.IsInf(dist[i], 1) {
				t.Errorf("l=%d: Vertex %d unreachable", l, i)
			}
		}
	}
}

func TestBMSSP_BoundaryConstraint(t *testing.T) {
	g := createLinearGraph(10)
	B := 5.0 // Small boundary
	algo := NewBMSSPAlgorithm(g, 2, B, []int{0})
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	// Check within/beyond the boundary behavior.
	for i := 0; i < 10; i++ {
		if float64(i) < B && dist[i] != float64(i) {
			t.Errorf("Vertex %d within boundary has wrong distance: got %f, want %f", i, dist[i], float64(i))
		}
	}
}

func TestBMSSP_LargeSparseGraph(t *testing.T) {
	n := 100
	g := createLinearGraph(n)
	S := []int{0, n / 2, n - 1}
	algo := NewBMSSPAlgorithm(g, 4, 1000.0, S)
	dist, err := algo.Solve()
	if err != nil {
		t.Fatalf("Solve() returned an error: %v", err)
	}

	for i := 0; i < n; i++ {
		if math.IsInf(dist[i], 1) {
			t.Errorf("Vertex %d has infinite distance", i)
		}
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
	g := &common.Graph{
		N:   4,
		Adj: make(map[int][]common.Edge),
	}
	edges := []common.Edge{
		{0, 1, 1.0},
		{1, 2, 2.0},
		{2, 3, 1.0},
		{3, 0, 5.0},
	}
	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges
	return g
}
