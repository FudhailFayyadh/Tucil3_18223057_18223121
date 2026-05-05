package puzzle

import "math"

type HeuristicFn func(s State, b *Board) int

func absi(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// H1: Manhattan distance to goal
func H1Manhattan(s State, b *Board) int {
	return absi(s.Row-b.GoalRow) + absi(s.Col-b.GoalCol)
}

// H2: Euclidean distance to goal (floor)
func H2Euclidean(s State, b *Board) int {
	dr := float64(s.Row - b.GoalRow)
	dc := float64(s.Col - b.GoalCol)
	return int(math.Sqrt(dr*dr + dc*dc))
}

// H3: Chebyshev distance to goal
func H3Chebyshev(s State, b *Board) int {
	dr := absi(s.Row - b.GoalRow)
	dc := absi(s.Col - b.GoalCol)
	if dr > dc {
		return dr
	}
	return dc
}

func GetHeuristic(name string) HeuristicFn {
	switch name {
	case "H2":
		return H2Euclidean
	case "H3":
		return H3Chebyshev
	default:
		return H1Manhattan
	}
}
