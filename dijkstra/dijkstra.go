package dijkstra

import (
	"container/heap"
	"errors"
	"math"
	"playground/common"
)

// DijkstraAlgorithm encapsulates the state for a run of Dijkstra's algorithm.
type DijkstraAlgorithm struct {
	graph    *common.Graph
	sources  []int
	boundary *float64 // A nil boundary means the search is unbounded.
}

// NewDijkstraAlgorithm creates a new solver for Dijkstra's algorithm.
func NewDijkstraAlgorithm(g *common.Graph, sources []int, boundary *float64) *DijkstraAlgorithm {
	return &DijkstraAlgorithm{
		graph:    g,
		sources:  sources,
		boundary: boundary,
	}
}

// Solve executes Dijkstra's algorithm based on the configured sources and boundary.
func (a *DijkstraAlgorithm) Solve() (map[int]float64, error) {
	if len(a.sources) == 0 {
		return nil, errors.New("dijkstra: at least one source vertex must be provided")
	}

	dist := make(map[int]float64)
	for i := 0; i < a.graph.N; i++ {
		dist[i] = math.Inf(1)
	}

	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)

	for _, s := range a.sources {
		if s >= 0 && s < a.graph.N {
			dist[s] = 0
			heap.Push(&pq, &common.DistEntry{Vertex: s, Dist: 0})
		}
	}

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*common.DistEntry)
		u := entry.Vertex
		d := entry.Dist

		if d > dist[u] {
			continue
		}
		if a.boundary != nil && dist[u] >= *a.boundary {
			continue
		}

		for _, edge := range a.graph.Adj[u] {
			v := edge.V
			newDist := dist[u] + edge.Weight

			if a.boundary != nil && newDist >= *a.boundary {
				continue
			}
			if newDist < dist[v] {
				dist[v] = newDist
				heap.Push(&pq, &common.DistEntry{Vertex: v, Dist: newDist})
			}
		}
	}

	return dist, nil
}
