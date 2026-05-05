package puzzle

import (
	"container/heap"
	"time"
)

type Node struct {
	State    State
	GCost    int
	Priority int
	Moves    []Direction
	States   []State
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Priority < pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Node)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := old[len(old)-1]
	*pq = old[:len(old)-1]
	return n
}

type IterLog struct {
	Iteration int
	State     State
	GCost     int
	Priority  int
	Moves     []Direction
}

type SolveResult struct {
	Found      bool
	Moves      []Direction
	States     []State
	TotalCost  int
	Iterations int
	TimeMs     float64
	Log        []IterLog
}

func solve(b *Board, priorityFn func(g int, s State) int) SolveResult {
	start := time.Now()
	init := b.InitialState()

	startMoves := []Direction{}
	startStates := []State{init}

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Node{
		State:    init,
		GCost:    0,
		Priority: priorityFn(0, init),
		Moves:    startMoves,
		States:   startStates,
	})

	// visited: state → best gCost seen
	visited := map[State]int{}
	iterations := 0
	var logs []IterLog

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*Node)
		iterations++

		logs = append(logs, IterLog{
			Iteration: iterations,
			State:     cur.State,
			GCost:     cur.GCost,
			Priority:  cur.Priority,
			Moves:     append([]Direction{}, cur.Moves...),
		})

		if best, seen := visited[cur.State]; seen && best <= cur.GCost {
			continue
		}
		visited[cur.State] = cur.GCost

		if b.IsGoal(cur.State) {
			elapsed := float64(time.Since(start).Microseconds()) / 1000.0
			return SolveResult{
				Found:      true,
				Moves:      cur.Moves,
				States:     cur.States,
				TotalCost:  cur.GCost,
				Iterations: iterations,
				TimeMs:     elapsed,
				Log:        logs,
			}
		}

		for _, dir := range AllDirections {
			res := Slide(b, cur.State, dir)
			if !res.Valid {
				continue
			}
			newG := cur.GCost + res.Cost
			if best, seen := visited[res.NewState]; seen && best <= newG {
				continue
			}
			newMoves := append(append([]Direction{}, cur.Moves...), dir)
			newStates := append(append([]State{}, cur.States...), res.NewState)
			heap.Push(pq, &Node{
				State:    res.NewState,
				GCost:    newG,
				Priority: priorityFn(newG, res.NewState),
				Moves:    newMoves,
				States:   newStates,
			})
		}
	}

	elapsed := float64(time.Since(start).Microseconds()) / 1000.0
	return SolveResult{Found: false, Iterations: iterations, TimeMs: elapsed, Log: logs}
}

func SolveUCS(b *Board) SolveResult {
	return solve(b, func(g int, s State) int { return g })
}

func SolveGBFS(b *Board, h HeuristicFn) SolveResult {
	return solve(b, func(g int, s State) int { return h(s, b) })
}

func SolveAStar(b *Board, h HeuristicFn) SolveResult {
	return solve(b, func(g int, s State) int { return g + h(s, b) })
}
