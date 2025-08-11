package bmssp

import (
	"container/heap"
	"playground/common"
)

type DataStructureD struct {
	pq     common.PriorityQueue
	inHeap map[int]bool
	M      int
}

func (d *DataStructureD) Initialize(M int) {
	d.M = M
	d.pq = make(common.PriorityQueue, 0)
	d.inHeap = make(map[int]bool)
}

func (d *DataStructureD) Insert(v int, dist float64) {
	if !d.inHeap[v] {
		heap.Push(&d.pq, &common.DistEntry{Vertex: v, Dist: dist})
		d.inHeap[v] = true
	}
}

func (d *DataStructureD) Pull() (int, []int) {
	if d.pq.Len() == 0 {
		return -1, nil
	}
	entry := heap.Pop(&d.pq).(*common.DistEntry)
	delete(d.inHeap, entry.Vertex)

	return int(entry.Dist), []int{entry.Vertex}
}

func (d *DataStructureD) BatchPrepend(entries []common.DistEntry) {
	for _, e := range entries {
		d.Insert(e.Vertex, e.Dist)
	}
}

func (d *DataStructureD) IsEmpty() bool {
	return d.pq.Len() == 0
}
