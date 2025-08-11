package bmssp

import (
	"container/heap"
	"math"
)

type Edge struct {
	U, V   int
	Weight float64
}

type Graph struct {
	N     int
	Edges []Edge
	Adj   map[int][]Edge
}

type DistEntry struct {
	vertex int
	dist   float64
	index  int
}

func BaseCase(B int, S []int, g *Graph, dist map[int]float64) (int, []int) {
	// Run Dijkstra's algorithm for base case
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	inHeap := make(map[int]bool)

	// Initialize with source vertices
	for _, s := range S {
		heap.Push(&pq, &DistEntry{vertex: s, dist: dist[s]})
		inHeap[s] = true
	}

	U := make([]int, 0)

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*DistEntry)
		u := entry.vertex
		delete(inHeap, u)

		if dist[u] >= float64(B) {
			break
		}

		U = append(U, u)

		// Relax edges
		for _, e := range g.Adj[u] {
			v := e.V
			newDist := dist[u] + e.Weight

			if newDist < dist[v] {
				dist[v] = newDist
				if !inHeap[v] && newDist < float64(B) {
					heap.Push(&pq, &DistEntry{vertex: v, dist: newDist})
					inHeap[v] = true
				}
			}
		}
	}

	// Compute actual boundary
	Bp := B
	for _, u := range U {
		if dist[u] < float64(Bp) {
			Bp = int(dist[u])
		}
	}

	return Bp, U
}

func FindPivots(B int, g *Graph, dist map[int]float64) ([]int, map[int]float64) {
	P := make([]int, 0)
	W := make(map[int]float64)

	for v := 0; v < g.N; v++ {
		if dist[v] < float64(B)/2 && dist[v] > 0 {
			P = append(P, v)
		}
		if dist[v] < float64(B) {
			W[v] = dist[v]
		}
	}

	return P, W
}

func BMSSP(l, B int, S []int, g *Graph, dist map[int]float64) (int, []int) {
	if dist == nil {
		dist = make(map[int]float64)
		for i := 0; i < g.N; i++ {
			dist[i] = math.Inf(1)
		}
		for _, s := range S {
			dist[s] = 0
		}
	}

	if l == 0 {
		Bp, U := BaseCase(B, S, g, dist)
		return Bp, U
	}

	P, W := FindPivots(B, g, dist)

	// If no pivots found, fall back to base case
	if len(P) == 0 {
		return BaseCase(B, S, g, dist)
	}

	var D DataStructureD
	M := 1 << uint(l-1) // 2^(l-1)
	D.Initialize(M, B)

	for _, x := range P {
		if d, exists := dist[x]; exists {
			D.Insert(x, d)
		}
	}

	U := make([]int, 0)
	k := 1 << uint(l) // k*2^l

	for len(U) < k && !D.IsEmpty() {
		Bi, Si := D.Pull()
		Bip, Ui := BMSSP(l-1, Bi, Si, g, dist)
		U = append(U, Ui...)

		K := make([]DistEntry, 0)

		for _, u := range Ui {
			for _, e := range g.Adj[u] {
				v := e.V
				newDist := dist[u] + e.Weight

				if newDist <= dist[v] {
					dist[v] = newDist

					if newDist >= float64(Bi) && newDist < float64(B) {
						D.Insert(v, newDist)
					} else if newDist <= float64(Bip) && newDist < float64(Bi) {
						K = append(K, DistEntry{vertex: v, dist: newDist})
					}
				}
			}
		}

		toAdd := K
		for _, s := range S {
			if dist[s] >= float64(Bip) && dist[s] < float64(Bi) {
				toAdd = append(toAdd, DistEntry{vertex: s, dist: dist[s]})
			}
		}
		D.BatchPrepend(toAdd)
	}

	// Compute final boundary B'
	Bp := B
	if len(U) > 0 {
		minBi := B
		for _, u := range U {
			if dist[u] < float64(minBi) {
				minBi = int(dist[u])
			}
		}
		if minBi < Bp {
			Bp = minBi
		}
	}

	// Add vertices from W with distance < B'
	for v, wDist := range W {
		if wDist < float64(Bp) {
			U = append(U, v)
		}
	}

	return Bp, U
}
