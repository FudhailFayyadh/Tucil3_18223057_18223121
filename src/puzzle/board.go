package puzzle

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	TilePath  = byte('*')
	TileWall  = byte('X')
	TileLava  = byte('L')
	TileGoal  = byte('O')
	TileStart = byte('Z')
)

type Board struct {
	N, M         int
	Tiles        [][]byte
	Costs        [][]int
	StartRow     int
	StartCol     int
	GoalRow      int
	GoalCol      int
	Checkpoints  []int // sorted unique digit values present on board
}

func ParseBoard(filename string) (*Board, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Read N M
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	parts := strings.Fields(scanner.Text())
	if len(parts) < 2 {
		return nil, fmt.Errorf("first line must contain N M")
	}
	n, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || n <= 0 || m <= 0 {
		return nil, fmt.Errorf("invalid N M values")
	}

	b := &Board{N: n, M: m, StartRow: -1, StartCol: -1, GoalRow: -1, GoalCol: -1}
	b.Tiles = make([][]byte, n)
	digitSet := map[int]bool{}

	// Read tile rows
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("expected %d tile rows, got %d", n, i)
		}
		line := scanner.Text()
		if len(line) != m {
			return nil, fmt.Errorf("row %d has length %d, expected %d", i, len(line), m)
		}
		b.Tiles[i] = []byte(line)
		for j := 0; j < m; j++ {
			ch := b.Tiles[i][j]
			switch ch {
			case TileStart:
				if b.StartRow >= 0 {
					return nil, fmt.Errorf("multiple start positions")
				}
				b.StartRow, b.StartCol = i, j
			case TileGoal:
				if b.GoalRow >= 0 {
					return nil, fmt.Errorf("multiple goal positions")
				}
				b.GoalRow, b.GoalCol = i, j
			case TilePath, TileWall, TileLava:
				// valid
			default:
				if ch >= '0' && ch <= '9' {
					digitSet[int(ch-'0')] = true
				} else {
					return nil, fmt.Errorf("unknown tile '%c' at (%d,%d)", ch, i, j)
				}
			}
		}
	}

	if b.StartRow < 0 {
		return nil, fmt.Errorf("no start position (Z) found")
	}
	if b.GoalRow < 0 {
		return nil, fmt.Errorf("no goal position (O) found")
	}

	// Read cost rows
	b.Costs = make([][]int, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("expected %d cost rows, got %d", n, i)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) != m {
			return nil, fmt.Errorf("cost row %d has %d values, expected %d", i, len(fields), m)
		}
		b.Costs[i] = make([]int, m)
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid cost at (%d,%d): %s", i, j, f)
			}
			b.Costs[i][j] = v
		}
	}

	// Build sorted checkpoints
	for d := range digitSet {
		b.Checkpoints = append(b.Checkpoints, d)
	}
	sort.Ints(b.Checkpoints)

	return b, nil
}

func (b *Board) IsGoal(s State) bool {
	return s.Row == b.GoalRow && s.Col == b.GoalCol && s.CheckpointIdx == len(b.Checkpoints)
}

func (b *Board) InitialState() State {
	return State{Row: b.StartRow, Col: b.StartCol, CheckpointIdx: 0}
}

// PrintWithActor renders the board with actor at (row,col).
// checkpointIdx is the index of the next required checkpoint; all digits
// with value < Checkpoints[checkpointIdx] have been collected and show as '*'.
func (b *Board) PrintWithActor(row, col, checkpointIdx int) string {
	var sb strings.Builder
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.M; j++ {
			tile := b.Tiles[i][j]
			if i == row && j == col {
				sb.WriteByte('Z')
			} else if tile == TileStart {
				sb.WriteByte('*')
			} else if tile >= '0' && tile <= '9' {
				digit := int(tile - '0')
				// if this digit has already been collected, show as '*'
				if checkpointIdx > 0 && (checkpointIdx >= len(b.Checkpoints) || digit < b.Checkpoints[checkpointIdx]) {
					sb.WriteByte('*')
				} else {
					sb.WriteByte(tile)
				}
			} else {
				sb.WriteByte(tile)
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
