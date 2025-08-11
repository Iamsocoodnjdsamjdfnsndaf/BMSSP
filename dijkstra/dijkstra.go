package dijkstra

import (
	"container/heap"
	"math"
	"playground/common"
)

func Dijkstra(g *common.Graph, s int) map[int]float64 {
	dist := make(map[int]float64)

	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	dist[s] = 0

	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &common.DistEntry{Vertex: s, Dist: 0})

	inHeap := make(map[int]bool)
	inHeap[s] = true

	processed := make(map[int]bool)

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*common.DistEntry)
		u := entry.Vertex
		delete(inHeap, u)

		if processed[u] {
			continue
		}
		processed[u] = true

		// Relax all edges from u
		for _, edge := range g.Adj[u] {
			v := edge.V
			newDist := dist[u] + edge.Weight

			if newDist < dist[v] {
				dist[v] = newDist
				if !processed[v] {
					heap.Push(&pq, &common.DistEntry{Vertex: v, Dist: newDist})
					inHeap[v] = true
				}
			}
		}
	}

	return dist
}

func DijkstraMultiSource(g *common.Graph, sources []int) map[int]float64 {
	dist := make(map[int]float64)

	// Initialize distances
	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	for _, s := range sources {
		dist[s] = 0
	}

	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)

	inHeap := make(map[int]bool)
	for _, s := range sources {
		heap.Push(&pq, &common.DistEntry{Vertex: s, Dist: 0})
		inHeap[s] = true
	}

	processed := make(map[int]bool)

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*common.DistEntry)
		u := entry.Vertex
		delete(inHeap, u)

		if processed[u] {
			continue
		}
		processed[u] = true

		for _, edge := range g.Adj[u] {
			v := edge.V
			newDist := dist[u] + edge.Weight

			if newDist < dist[v] {
				dist[v] = newDist
				if !processed[v] {
					heap.Push(&pq, &common.DistEntry{Vertex: v, Dist: newDist})
					inHeap[v] = true
				}
			}
		}
	}

	return dist
}

func DijkstraBounded(g *common.Graph, sources []int, B float64) map[int]float64 {
	dist := make(map[int]float64)

	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	for _, s := range sources {
		dist[s] = 0
	}

	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)

	inHeap := make(map[int]bool)
	for _, s := range sources {
		heap.Push(&pq, &common.DistEntry{Vertex: s, Dist: 0})
		inHeap[s] = true
	}

	processed := make(map[int]bool)

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*common.DistEntry)
		u := entry.Vertex

		if dist[u] >= B {
			break
		}

		delete(inHeap, u)

		if processed[u] {
			continue
		}
		processed[u] = true

		for _, edge := range g.Adj[u] {
			v := edge.V
			newDist := dist[u] + edge.Weight

			if newDist < dist[v] && newDist < B {
				dist[v] = newDist
				if !processed[v] {
					heap.Push(&pq, &common.DistEntry{Vertex: v, Dist: newDist})
					inHeap[v] = true
				}
			}
		}
	}

	return dist
}
