package bmssp

import (
	"container/heap"
	"math"
	"playground/common"
	"slices"
)

type BMSSPAlgorithm struct {
	graph *common.Graph
	dist  map[int]float64
	l     int
	B     float64
	S     []int
}

func NewBMSSPAlgorithm(g *common.Graph, l int, B float64, S []int) *BMSSPAlgorithm {
	dist := make(map[int]float64, g.N)
	for i := 0; i < g.N; i++ {
		dist[i] = math.Inf(1)
	}
	for _, s := range S {
		dist[s] = 0
	}
	return &BMSSPAlgorithm{
		graph: g,
		dist:  dist,
		l:     l,
		B:     B,
		S:     S,
	}
}

func (a *BMSSPAlgorithm) Solve() (map[int]float64, error) {
	a.bmsspRecursive(a.l, a.B, a.S)
	return a.dist, nil
}

// k = floor(log(n)^(1/3)), t = floor(log(n)^(2/3)), each ≥ 1
func (a *BMSSPAlgorithm) kt() (int, int) {
	n := float64(a.graph.N)
	if n < 2 {
		n = 2
	}
	ln := math.Log(n)
	k := int(math.Floor(math.Pow(ln, 1.0/3.0)))
	t := int(math.Floor(math.Pow(ln, 2.0/3.0)))
	if k < 1 {
		k = 1
	}
	if t < 1 {
		t = 1
	}
	return k, t
}

func (a *BMSSPAlgorithm) bmsspRecursive(l int, B float64, S []int) (float64, []int) {
	if l == 0 {
		return a.baseCaseSingletonOrSplit(B, S)
	}

	P, W := a.findPivots(B, S)

	// No pivots ⇒ successful execution: B' = B, add W' = { x in W : d̂[x] < B }.
	if len(P) == 0 {
		U := make([]int, 0, len(W))
		seen := make(map[int]bool, len(W))
		for v, dv := range W {
			if dv < B && !seen[v] {
				seen[v] = true
				U = append(U, v)
			}
		}
		slices.Sort(U)
		return B, U
	}

	k, t := a.kt()
	M := 1 << uint((l-1)*t)
	var D DataStructureD
	D.Initialize(M, B)

	for _, x := range P {
		if d, exists := a.dist[x]; exists {
			D.Insert(x, d)
		}
	}

	threshold := k * (1 << uint(l*t))
	U := make([]int, 0, threshold)
	seenU := make(map[int]bool)
	lastBip := B

	for len(U) < threshold && !D.IsEmpty() {
		Bi, Si := D.Pull()
		Bip, Ui := a.bmsspRecursive(l-1, Bi, Si)
		lastBip = Bip

		// Dedup Ui before adding to U
		for _, u := range Ui {
			if !seenU[u] {
				seenU[u] = true
				U = append(U, u)
			}
		}

		// After recursion, do <= relaxations but NEVER write distances ≥ B
		K := make([]common.DistEntry, 0)
		for _, u := range Ui {
			for _, e := range a.graph.Adj[u] {
				v := e.V
				newDist := a.dist[u] + e.Weight
				if newDist < B && newDist <= a.dist[v] {
					a.dist[v] = newDist
					if newDist >= Bi { // Bi ≤ newDist < B
						D.Insert(v, newDist)
					} else if newDist >= Bip { // Bip ≤ newDist < Bi
						K = append(K, common.DistEntry{Vertex: v, Dist: newDist})
					}
				}
			}
		}

		// Also include sources from Si whose key now falls in [Bip, Bi)
		for _, sNode := range Si {
			if d, ok := a.dist[sNode]; ok && d < B && d >= Bip && d < Bi {
				K = append(K, common.DistEntry{Vertex: sNode, Dist: d})
			}
		}

		D.BatchPrepend(K)
	}

	Bp := math.Min(lastBip, B)

	// Add W' = { x in W : d̂[x] < B' } (dedup against U)
	for v, dv := range W {
		if dv < Bp && !seenU[v] {
			seenU[v] = true
			U = append(U, v)
		}
	}

	slices.Sort(U)
	return Bp, U
}

// Base case: enforce singleton, split otherwise
func (a *BMSSPAlgorithm) baseCaseSingletonOrSplit(B float64, S []int) (float64, []int) {
	if len(S) == 1 {
		return a.baseCase(B, S[0])
	}
	minBp := B
	seen := map[int]bool{}
	U := []int{}
	for _, s := range S {
		Bpi, Ui := a.baseCase(B, s)
		if Bpi < minBp {
			minBp = Bpi
		}
		for _, u := range Ui {
			if !seen[u] {
				seen[u] = true
				U = append(U, u)
			}
		}
	}
	slices.Sort(U)
	return minBp, U
}

// Robust single-source bounded Dijkstra (lazy decrease-key) that NEVER writes dist ≥ B
func (a *BMSSPAlgorithm) baseCase(B float64, s int) (float64, []int) {
	pq := make(common.PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &common.DistEntry{Vertex: s, Dist: a.dist[s]})

	U := make([]int, 0)
	seenU := make(map[int]bool)

	for pq.Len() > 0 {
		entry := heap.Pop(&pq).(*common.DistEntry)
		u := entry.Vertex

		// Skip stale entries
		if entry.Dist != a.dist[u] {
			continue
		}
		// If the smallest key is ≥ B, we're at the boundary
		if entry.Dist >= B {
			heap.Push(&pq, entry) // put back so pq[0] is the boundary key
			break
		}

		if !seenU[u] {
			seenU[u] = true
			U = append(U, u)
		}

		for _, e := range a.graph.Adj[u] {
			v := e.V
			newDist := entry.Dist + e.Weight
			// Only write & push if strictly below B
			if newDist < B && newDist <= a.dist[v] {
				a.dist[v] = newDist
				heap.Push(&pq, &common.DistEntry{Vertex: v, Dist: newDist})
			}
		}
	}

	Bp := B
	if pq.Len() > 0 {
		minBoundaryDist := pq[0].Dist
		if minBoundaryDist < Bp {
			Bp = minBoundaryDist
		}
	}
	return Bp, U
}

// Pivot-finding (k-step <=-relaxations bounded by B, equality forest, roots with size ≥ k)
func (a *BMSSPAlgorithm) findPivots(B float64, S []int) ([]int, map[int]float64) {
	kParam, _ := a.kt()

	WMap := make(map[int]bool, len(S))
	for _, sNode := range S {
		WMap[sNode] = true
	}

	pred := make(map[int]int)

	// BFS-like k rounds of <= relaxations within bound B
	WFrontiers := make([][]int, kParam+1)
	WFrontiers[0] = S
	for i := 1; i <= kParam; i++ {
		WFrontiers[i] = make([]int, 0)
		frontierMap := make(map[int]bool)
		for _, u := range WFrontiers[i-1] {
			for _, e := range a.graph.Adj[u] {
				v := e.V
				newDist := a.dist[u] + e.Weight
				if newDist < B && newDist <= a.dist[v] {
					a.dist[v] = newDist
					pred[v] = u
					WMap[v] = true
					if !frontierMap[v] {
						WFrontiers[i] = append(WFrontiers[i], v)
						frontierMap[v] = true
					}
				}
			}
		}
	}

	WDistMap := make(map[int]float64, len(WMap))
	for v := range WMap {
		WDistMap[v] = a.dist[v]
	}

	// If W is too large, choose all S as pivots
	if len(WDistMap) > kParam*len(S) {
		return S, WDistMap
	}

	// Build equality-forest and compute subtree sizes
	children := make(map[int][]int)
	for v, u := range pred {
		if WMap[v] && WMap[u] {
			children[u] = append(children[u], v)
		}
	}

	treeSizes := make(map[int]int)
	var dfs func(u int) int
	dfs = func(u int) int {
		if size, ok := treeSizes[u]; ok {
			return size
		}
		size := 1
		for _, vv := range children[u] {
			size += dfs(vv)
		}
		treeSizes[u] = size
		return size
	}

	P := make([]int, 0)
	for _, u := range S {
		if _, ok := treeSizes[u]; !ok {
			dfs(u)
		}
		if treeSizes[u] >= kParam {
			P = append(P, u)
		}
	}

	return P, WDistMap
}
