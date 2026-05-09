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

func (prioQ PriorityQueue) Len() int {
	return len(prioQ)
} 
func (prioQ PriorityQueue) Less(i, j int) bool {
	return prioQ[i].Priority < prioQ[j].Priority
}

func (prioQ PriorityQueue) Swap(i, j int) {
	temp := prioQ [i]
	prioQ[i] = prioQ[j]
	prioQ[j] = temp
}

func (prioQ PriorityQueue) Push(x interface{}) {
	newNode := x.(*Node)
	*prioQ = append(*prioQ, newNode)
}

func (prioQ *PriorityQueue) Pop() interface{} {
	prioQ_temp := *prioQ
	length := len(prioQ_temp)

	// pop disini
	result := prioQ[length-1]
	*prioQ = prioQ_temp[:length-1]

	return result
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

func solve(board *Board, prioCount func(g int, s State) int) SolveResult {
	startTime := time.Now()

	awal := board.InitialState()

	nodeAwal := &Node{
		state: awal,
		GCost: 0,
		Priority: prioCount(0, awal),
		Moves: []Direction{},
		States: []State{awal},
	}

	prioQ := &PriorityQueue{}
	heap.Init(prioQ)
	heap.Push(prioQ, nodeAwal)
	
	visitedList := map[State]int{}

	cntIter :=0
	var iterLog []IterLog

	for prioQ.Len() > 0 {
		node := heap.Pop(prioQ).(*Node)
		cntIter++

		iterLog = append(iterLog, IterLog{
			Iteration: cntIter,
			State: node.State,
			GCost: node.GCost,
			Priority: node.Priority,
			Moves: append([]Direction{}, node.Moves...),
		})

		bestCost, visitedBLock := visitedList[node.State]

		if (visitedBLock && bestCost <= node.GCost) {
			continue
		}

		visitedList[node.State] = node.GCost

		if board.IsGoal(node.State) {
			finishTime := float64(time.Since(startTime).Microseconds()) / 1000.0

			return result {
				Found: true,
				Moves: node.Moves,
				States: node.States,
				TotalCost: node.GCost,
				Iterations: cntIter,
				TimeMs: finishTime,
				Log: iterLog
			}
		}

		for _, direction := range AllDirections {
			move := Slide(board, node.State, direction)

			if !move.Valid {
				continue
			}

			newCost := node.GCost + move.Cost

			newBestCost, newVisitedBlock := visitedList[move.NewState]
			
			if (newVisitedBlock && newBestCost <= newCost) {
				continue
			}

			newStep := append([]Direction{}, node.Moves...)
			newStep = append(newStep, direction)
			
			newState := append([]State{}, node.States...)
			newState = append(newState, move.NewState)

			newNode := &Node{
				State: move.NewState,
				GCost: newCost,
				Priority: prioCount(newCost, move.NewState),
				Moves: newStep,
				States: newState,
			}
			heap.Push(prioQ, newNode)
		}
	}

	finishTime := float64(time.Since(startTime).Microseconds()) / 1000.0
	return result {
		Found: false,
		Iterations: cntIter,
		TimeMs: finishTime,
		Log: iterLog
	}
}

func SolveUCS(board *Board) SolveResult {
	prioCount := func(g int, s State) int {
		return g
	}
	return solve(board, prioCount)
}

func SolveGBFS(board *Board, heuristik HeuristicFn) SolveResult {
	prioCount := func(g int, s State) int {
		return heuristik(s, board)
	}
	return solve(board, prioCount)
}

func SolveAStar(board *Board, heuristik HeuristicFn) SolveResult {
	prioCount := func(g int, s State) int {
		return g + heuristik (s, board)
	}
	return solve(board, prioCount)
}
