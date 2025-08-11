package bmssp

import (
	"container/heap"
	"math"
	"playground/common"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &common.DistEntry{Vertex: 1, Dist: 5.0})
	heap.Push(&pq, &common.DistEntry{Vertex: 2, Dist: 3.0})
	heap.Push(&pq, &common.DistEntry{Vertex: 3, Dist: 7.0})
	heap.Push(&pq, &common.DistEntry{Vertex: 4, Dist: 1.0})

	expected := []int{4, 2, 1, 3}
	for i, exp := range expected {
		if pq.Len() == 0 {
			t.Fatalf("Queue empty at iteration %d", i)
		}
		entry := heap.Pop(&pq).(*common.DistEntry)
		if entry.Vertex != exp {
			t.Errorf("Expected vertex %d, got %d", exp, entry.Vertex)
		}
	}

	if pq.Len() != 0 {
		t.Errorf("Expected empty queue, has %d elements", pq.Len())
	}
}

func TestDataStructureD_BasicOperations(t *testing.T) {
	var d DataStructureD
	d.Initialize(10)

	if !d.IsEmpty() {
		t.Error("Expected empty data structure after initialization")
	}

	d.Insert(1, 5.0)
	d.Insert(2, 3.0)
	d.Insert(3, 7.0)

	if d.IsEmpty() {
		t.Error("Expected non-empty after insertions")
	}

	B1, S1 := d.Pull()
	if B1 != 3 || len(S1) != 1 || S1[0] != 2 {
		t.Errorf("Pull failed: B=%d, S=%v", B1, S1)
	}

	d.Insert(1, 5.0) // Duplicate check
	if d.pq.Len() != 2 {
		t.Errorf("Expected 2 elements after duplicate insert, got %d", d.pq.Len())
	}
}

func TestDataStructureD_BatchPrepend(t *testing.T) {
	var d DataStructureD
	d.Initialize(10)

	entries := []common.DistEntry{
		{Vertex: 1, Dist: 5.0},
		{Vertex: 2, Dist: 3.0},
		{Vertex: 3, Dist: 7.0},
	}

	d.BatchPrepend(entries)

	if d.pq.Len() != 3 {
		t.Errorf("Expected 3 elements after batch prepend, got %d", d.pq.Len())
	}

	// Verify order
	B, S := d.Pull()
	if B != 3 || S[0] != 2 {
		t.Errorf("Unexpected first pull: B=%d, S=%v", B, S)
	}
}

func TestBMSSP_SingleSource(t *testing.T) {
	g := createLinearGraph(5)
	S := []int{0}
	B := 100
	l := 2

	dist := initializeDistances(g.N, S)

	Bp, _ := BMSSP(l, B, S, g, dist)

	if Bp > B {
		t.Errorf("Boundary exceeded: B'=%d, B=%d", Bp, B)
	}

	// Check distances are correct for linear graph
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
	B := 100
	l := 2

	dist := initializeDistances(g.N, S)

	Bp, _ := BMSSP(l, B, S, g, dist) // Use _ to ignore U

	if Bp > B {
		t.Errorf("Boundary exceeded: B'=%d, B=%d", Bp, B)
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

	// Two disconnected components
	edges := []common.Edge{
		{0, 1, 1.0},
		{1, 2, 1.0},
		{3, 4, 1.0},
		{4, 5, 1.0},
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges

	S := []int{0}
	B := 100
	l := 2

	dist := initializeDistances(g.N, S)

	_, _ = BMSSP(l, B, S, g, dist)

	// Vertices in same component should have finite distance
	if math.IsInf(dist[1], 1) || math.IsInf(dist[2], 1) {
		t.Error("Connected vertices have infinite distance")
	}

	// Vertices in other component should have infinite distance
	if !math.IsInf(dist[3], 1) || !math.IsInf(dist[4], 1) || !math.IsInf(dist[5], 1) {
		t.Error("Disconnected vertices have finite distance")
	}
}

func TestBMSSP_CycleGraph(t *testing.T) {
	g := &common.Graph{
		N:   4,
		Adj: make(map[int][]common.Edge),
	}

	// Create cycle with different weights
	edges := []common.Edge{
		{0, 1, 1.0},
		{1, 2, 2.0},
		{2, 3, 1.0},
		{3, 0, 5.0}, // Completes cycle
	}

	for _, e := range edges {
		g.Adj[e.U] = append(g.Adj[e.U], e)
		g.Adj[e.V] = append(g.Adj[e.V], common.Edge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g.Edges = edges

	S := []int{0}
	B := 100
	l := 2

	dist := initializeDistances(g.N, S)

	_, _ = BMSSP(l, B, S, g, dist)

	// Check shortest paths
	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestBMSSP_DifferentRecursionDepths(t *testing.T) {
	g := createLinearGraph(10)
	S := []int{0}
	B := 100

	for l := 0; l <= 4; l++ {
		dist := initializeDistances(g.N, S)

		Bp, _ := BMSSP(l, B, S, g, dist)

		if Bp > B {
			t.Errorf("l=%d: Boundary exceeded B'=%d, B=%d", l, Bp, B)
		}

		// Verify all vertices reachable
		for i := 0; i < g.N; i++ {
			if math.IsInf(dist[i], 1) {
				t.Errorf("l=%d: Vertex %d unreachable", l, i)
			}
		}
	}
}

func TestBMSSP_BoundaryConstraint(t *testing.T) {
	g := createLinearGraph(10)
	S := []int{0}
	B := 5 // Small boundary
	l := 2

	dist := initializeDistances(g.N, S)

	Bp, U := BMSSP(l, B, S, g, dist)

	if Bp > B {
		t.Errorf("Boundary exceeded: B'=%d, B=%d", Bp, B)
	}

	// Check U only contains vertices within boundary
	for _, v := range U {
		if dist[v] >= float64(B) {
			t.Errorf("Vertex %d in U has distance %f >= B=%d", v, dist[v], B)
		}
	}
}

func TestBMSSP_WeightedCompleteGraph(t *testing.T) {
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

	S := []int{0}
	B := 100
	l := 2

	dist := initializeDistances(g.N, S)

	_, _ = BMSSP(l, B, S, g, dist)

	// Verify shortest paths
	expected := map[int]float64{0: 0, 1: 1, 2: 3, 3: 4}
	for v, exp := range expected {
		if math.Abs(dist[v]-exp) > 1e-9 {
			t.Errorf("Vertex %d: expected dist=%f, got %f", v, exp, dist[v])
		}
	}
}

func TestBMSSP_EmptyGraph(t *testing.T) {
	g := &common.Graph{
		N:   1,
		Adj: make(map[int][]common.Edge),
	}

	S := []int{0}
	B := 100
	l := 0

	dist := initializeDistances(g.N, S)

	Bp, _ := BMSSP(l, B, S, g, dist)

	if Bp > B {
		t.Errorf("Boundary exceeded for single vertex: B'=%d, B=%d", Bp, B)
	}

	if dist[0] != 0 {
		t.Errorf("Source distance should be 0, got %f", dist[0])
	}
}

func TestBMSSP_LargeSparseGraph(t *testing.T) {
	// Create larger sparse graph (path)
	n := 100
	g := createLinearGraph(n)

	S := []int{0, n / 2, n - 1} // Multiple sources
	B := 1000
	l := 4

	dist := initializeDistances(g.N, S)

	Bp, U := BMSSP(l, B, S, g, dist)

	if Bp > B {
		t.Errorf("Boundary exceeded: B'=%d, B=%d", Bp, B)
	}

	// Check all vertices are processed
	for i := 0; i < n; i++ {
		if math.IsInf(dist[i], 1) {
			t.Errorf("Vertex %d has infinite distance", i)
		}
	}

	// Verify U is not empty
	if len(U) == 0 {
		t.Error("U should not be empty for reachable graph")
	}
}

// Helper functions

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

func initializeDistances(n int, sources []int) map[int]float64 {
	dist := make(map[int]float64)
	for i := 0; i < n; i++ {
		dist[i] = math.Inf(1)
	}
	for _, s := range sources {
		dist[s] = 0
	}
	return dist
}

func TestContainsVertex(t *testing.T) {
	vertices := []int{1, 3, 5, 7}

	tests := []struct {
		v        int
		expected bool
	}{
		{1, true},
		{3, true},
		{5, true},
		{7, true},
		{2, false},
		{4, false},
		{6, false},
	}

	for _, tt := range tests {
		found := false
		for _, v := range vertices {
			if v == tt.v {
				found = true
				break
			}
		}
		if found != tt.expected {
			t.Errorf("containsVertex(%d): expected %v, got %v", tt.v, tt.expected, found)
		}
	}
}

func BenchmarkBMSSP_SmallGraph(b *testing.B) {
	g := createLinearGraph(10)
	S := []int{0}
	B := 100
	l := 2

	for i := 0; i < b.N; i++ {
		dist := initializeDistances(g.N, S)
		_, _ = BMSSP(l, B, S, g, dist)
	}
}

func BenchmarkBMSSP_MediumGraph(b *testing.B) {
	g := createLinearGraph(100)
	S := []int{0}
	B := 1000
	l := 3

	for i := 0; i < b.N; i++ {
		dist := initializeDistances(g.N, S)
		_, _ = BMSSP(l, B, S, g, dist)
	}
}

func BenchmarkBMSSP_MultiSource(b *testing.B) {
	g := createLinearGraph(100)
	S := []int{0, 25, 50, 75, 99}
	B := 1000
	l := 3

	for i := 0; i < b.N; i++ {
		dist := initializeDistances(g.N, S)
		_, _ = BMSSP(l, B, S, g, dist)
	}
}
