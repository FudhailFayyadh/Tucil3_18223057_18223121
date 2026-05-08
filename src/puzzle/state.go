package puzzle

type State struct {
	Row, Col      int
	CheckpointIdx int
}

type Direction byte

const (
	DirUp    Direction = 'U'
	DirDown  Direction = 'D'
	DirLeft  Direction = 'L'
	DirRight Direction = 'R'
)

var AllDirections = []Direction{DirUp, DirDown, DirLeft, DirRight}

type SlideResult struct {
	NewState State
	Cost     int
	Valid    bool
}

func dirDelta(d Direction) (int, int) {
	switch d {
	case DirUp:
		return -1, 0
	case DirDown:
		return 1, 0
	case DirLeft:
		return 0, -1
	default:
		return 0, 1
	}
}

func Slide(b *Board, s State, dir Direction) SlideResult {
	dr, dc := dirDelta(dir)
	rr, cc := s.Row, s.Col
	cost := 0
	idx := s.CheckpointIdx

	for {
		nr, nc := rr+dr, cc+dc

		// out of bounds → game over
		if nr < 0 || nr >= b.N || nc < 0 || nc >= b.M {
			return SlideResult{Valid: false}
		}

		tile := b.Tiles[nr][nc]

		// wall → stop before it
		if tile == TileWall {
			break
		}

		// move onto next tile
		rr, cc = nr, nc
		cost += b.Costs[rr][cc]

		if tile == TileLava {
			return SlideResult{Valid: false}
		}

		if tile >= '0' && tile <= '9' {
			n := int(tile - '0')
			if idx < len(b.Checkpoints) {
				nextReq := b.Checkpoints[idx]
				if n == nextReq {
					idx++
				} else if n > nextReq {
					// out-of-order checkpoint
					return SlideResult{Valid: false}
				}
				// n < nextReq → already collected, treat as normal
			}
			// idx == len means all done, digit is normal now
		}

	}

	if rr == s.Row && cc == s.Col {
		return SlideResult{Valid: false}
	}

	return SlideResult{
		NewState: State{Row: rr, Col: cc, CheckpointIdx: idx},
		Cost:     cost,
		Valid:    true,
	}
}
