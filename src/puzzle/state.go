package puzzle

type State struct {
	Row, Col int
	CheckpointIdx int
}

type Direction byte

const (
	DirUp Direction = 'U'
	DirDown Direction = 'D'
	DirLeft Direction = 'L'
	DirRight Direction = 'R'
)

var AllDirections = []Direction{DirUp, DirDown, DirLeft, DirRight}

type SlideResult struct {
	NewState State
	Cost int
	Valid bool
}

func dirDelta(dir Direction) (int, int) {
	var diffRow int
	var diffCol int

	if dir == DirUp {
		diffRow = -1
		diffCol = 0
	} else if dir == DirDown {
		diffRow = 1
		diffCol = 0
	} else if dir == DirLeft {
		diffRow = 0
		diffCol = -1
	} else {
		diffRow = 0
		diffCol = 1
	}

	return diffRow, diffCol
}

func Slide(board *Board, state State, dir Direction) SlideResult {
	diffRow, diffCol := dirDelta(dir)

	curRow := state.Row
	curCol := state.Col
	checkpointIdx := state.CheckpointIdx
	totalCost := 0

	for {
		nextRow := curRow + diffRow
		nextCol := curCol + diffCol
		
		// cek batas map dlu
		if nextRow < 0 || nextRow >= board.N || nextCol < 0 || nextCol >= board.M {
			var failResult SlideResult
			
			failResult.Valid = false
			return failResult
		}

		// cek tiles nya
		nextTile := board.Tiles[nextRow][nextCol]

		if nextTile == TileWall {
			break // Kalau tembok, berhenti biar ga error
		}

		curRow = nextRow
		curCol = nextCol

		stepCost := board.Costs[curRow][curCol]
		totalCost = totalCost + stepCost

		if nextTile == TileLava { // kena lava, dead
			var deathResult SlideResult
			
			deathResult.Valid = false
			return deathResult
		}

		if nextTile >= '0' && nextTile <= '9' {
			tileDigit := int (nextTile - '0')

			if checkpointIdx < len(board.Checkpoints) {
				reqDigit := board.Checkpoints[checkpointIdx]

				if tileDigit == reqDigit {
					checkpointIdx = checkpointIdx + 1

				} else if tileDigit > reqDigit { // salah jalan
					var wrongJumpResult SlideResult

					wrongJumpResult.Valid = false
					return wrongJumpResult
				}

			}

		}

	}

	if curRow == state.Row && curCol == state.Col { // ga gerak
		var stillResult SlideResult

		stillResult.Valid = false
		return stillResult
	}

	var resultState State
	resultState.Row = curRow
	resultState.Col = curCol
	resultState.CheckpointIdx = checkpointIdx

	var slideResult SlideResult
	slideResult.Valid = true
	slideResult.Cost = totalCost
	slideResult.NewState = resultState
	
	return slideResult
}
