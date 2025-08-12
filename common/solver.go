package common

type ShortestPathSolver interface {
	// Solve executes the algorithm and returns the final distance map.
	Solve() (map[int]float64, error)
}
