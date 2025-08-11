package bmssp

import "container/heap"

type DataStructureD struct {
	pq     PriorityQueue
	inHeap map[int]bool
	M      int
}

func (d *DataStructureD) Initialize(M, B int) {
	d.M = M
	d.pq = make(PriorityQueue, 0)
	d.inHeap = make(map[int]bool)
}

func (d *DataStructureD) Insert(v int, dist float64) {
	if !d.inHeap[v] {
		heap.Push(&d.pq, &DistEntry{vertex: v, dist: dist})
		d.inHeap[v] = true
	}
}

func (d *DataStructureD) Pull() (int, []int) {
	if d.pq.Len() == 0 {
		return -1, nil
	}
	entry := heap.Pop(&d.pq).(*DistEntry)
	delete(d.inHeap, entry.vertex)

	// Pull returns Bi and Si (boundary and vertex set)
	return int(entry.dist), []int{entry.vertex}
}

func (d *DataStructureD) BatchPrepend(entries []DistEntry) {
	for _, e := range entries {
		d.Insert(e.vertex, e.dist)
	}
}

func (d *DataStructureD) IsEmpty() bool {
	return d.pq.Len() == 0
}
