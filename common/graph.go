package common

type Graph struct {
	N     int
	Edges []Edge
	Adj   map[int][]Edge
}
