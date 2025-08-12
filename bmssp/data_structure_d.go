package bmssp

import (
	"container/heap"
	"playground/common"
)

// DataStructureD is a specialised priority queue for the BMSSP algorithm.
type DataStructureD struct {
	pq     common.PriorityQueue
	inHeap map[int]*common.DistEntry
	M      int
	B      float64
}

// Initialize sets up the DataStructureD with a batch size M and upper bound B.
func (d *DataStructureD) Initialize(M int, B float64) {
	d.M = M
	d.B = B
	d.pq = make(common.PriorityQueue, 0)
	d.inHeap = make(map[int]*common.DistEntry)
	heap.Init(&d.pq)
}

// Insert adds a vertex with its distance, updating the value if a shorter path is found.
func (d *DataStructureD) Insert(v int, dist float64) {
	if dist >= d.B {
		return
	}
	if entry, ok := d.inHeap[v]; ok {
		if dist < entry.Dist {
			entry.Dist = dist
			heap.Fix(&d.pq, entry.Index)
		}
		return
	}
	item := &common.DistEntry{Vertex: v, Dist: dist}
	heap.Push(&d.pq, item)
	d.inHeap[v] = item
}

// Pull pops up to M items, but NEVER splits ties: if the next key equals the last
// popped key, we keep popping to drain the entire tie group. Returns (Bi, S')
// where Bi is the next strictly larger key (or B if none).
func (d *DataStructureD) Pull() (float64, []int) {
	if d.pq.Len() == 0 {
		return d.B, nil
	}

	Sprime := make([]int, 0, d.M)

	// Pop the first item to determine the tie key.
	first := heap.Pop(&d.pq).(*common.DistEntry)
	delete(d.inHeap, first.Vertex)
	Sprime = append(Sprime, first.Vertex)
	popped := 1
	tieKey := first.Dist

	for d.pq.Len() > 0 {
		next := d.pq[0]
		if popped >= d.M && next.Dist > tieKey {
			break
		}
		entry := heap.Pop(&d.pq).(*common.DistEntry)
		delete(d.inHeap, entry.Vertex)
		Sprime = append(Sprime, entry.Vertex)
		popped++
	}

	// Bi = next strictly larger key, capped at B; if none, Bi = B.
	if d.pq.Len() == 0 {
		return d.B, Sprime
	}
	bi := d.pq[0].Dist
	if bi > d.B {
		bi = d.B
	}
	return bi, Sprime
}

// BatchPrepend inserts a list of entries, keeping the smallest distance per vertex.
func (d *DataStructureD) BatchPrepend(entries []common.DistEntry) {
	if len(entries) == 0 {
		return
	}
	best := make(map[int]float64, len(entries))
	for _, e := range entries {
		if cur, ok := best[e.Vertex]; !ok || e.Dist < cur {
			best[e.Vertex] = e.Dist
		}
	}
	for v, dist := range best {
		d.Insert(v, dist)
	}
}

// IsEmpty checks if the data structure is empty.
func (d *DataStructureD) IsEmpty() bool {
	return d.pq.Len() == 0
}
