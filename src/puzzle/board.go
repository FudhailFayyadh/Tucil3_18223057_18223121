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
	TilePath = byte('*')
	TileWall = byte('X')
	TileLava = byte('L')
	TileGoal = byte('O')
	TileStart = byte('Z')
)

type Board struct {
	N, M int
	Tiles [][]byte
	Costs [][]int
	StartRow int
	StartCol int
	GoalRow int
	GoalCol int
	Checkpoints []int 
}

func ParseBoard(namaFile string) (*Board, error) {
	file, err:= os.Open(namaFile)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka fila: %v", err)
		
	}
	defer file.Close()

	read := bufio.NewScanner(file)
	if !read.Scan() {
		return nil, fmt.Errorf("File kosong")
	}

	rowOne := read.Text()
	textRow := strings.Fields(rowOne)

	if len(textRow) <2 {
		return nil, fmt.Errorf("Harus memberikan 2 angka")
	}

	nBaris, err1 := strconv.Atoi(textRow[0])
	nKolom, err2 := strconv.Atoi(textRow[1])

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("Angka yang diberikan tidak valid")
	}

	board := &Board {
		N: nBaris,
		M: nKolom,
		StartRow : -1,
		StartCol : -1,
		GoalRow : -1,
		GoalCol : -1,
	}

	board.Tiles = make([][]byte, nBaris)
	daftarAngka := map[int]bool{}

	for i:=0; i<nBaris; i++ {
		if !read.Scan() {
			return nil, fmt.Errorf("Baris ke-%d kosong", i)
		}

		textBaris := read.Text()
		if len(textBaris) != nKolom {
			return nil, fmt.Errorf("Panjang baris tidak sesuai")
		}

		board.Tiles[i] = []byte(textBaris)
		for j:=0; j<nKolom; j++ {
			karakater := board.Tiles[i][j]

			if karakater == 'Z' { // tile start
				if board.StartRow != -1 {
					return nil, fmt.Errorf("Start tidak boleh lebih dari 1")
				}
				board.StartRow, board.StartCol = i, j
			} else if karakater == 'O' { // tile goal
				if board.GoalRow != -1 {
					return nil, fmt.Errorf("Goal tidak boleh lebih dari 1")
				}
				board.GoalRow, board.GoalCol = i, j
			} else if karakater >= '0' && karakater <= '9' {
				angka := int(karakater - '0')
				daftarAngka[angka] = true
			}
		}
	}

	if board.StartRow == -1 || board.GoalRow == -1 {
		return nil, fmt.Errorf("Start / Goal tidak ditemukan")
	}

	board.Costs = make ([][]int, nBaris)

	for i:=0; i < nBaris; i++ {
		if !read.Scan() {
			return nil, fmt.Errorf("Baris ke-%d kosong", i)
		}

		textAngka := strings.Fields(read.Text())
		if len(textAngka) != nKolom {
			return nil, fmt.Errorf("Panjang baris tidak sesuai")
		}

		board.Costs[i] = make([]int, nKolom)
		for j:=0; j<nKolom; j++ {
			biaya, err :=strconv.Atoi(textAngka[j])

			if err != nil {
				return nil, fmt.Errorf("...")
			}
			board.Costs[i][j] = biaya
		}
	}
	for angka := range daftarAngka {
		board.Checkpoints = append(board.Checkpoints, angka)
	}
	sort.Ints(board.Checkpoints)

	for i, cp := range board.Checkpoints {
		if cp != i {
			return nil, fmt.Errorf("sequence checkpoint tidak valid: harus 0,1,2,...,n tanpa skip")
		}
	}

	return board, nil
}

func (board *Board) IsGoal(status State) bool {
	if status.Row == board.GoalRow {
		if status.Col == board.GoalCol {
			total := len(board.Checkpoints)

			if status.CheckpointIdx == total {
				return true
			}
		}
	}
	return false 
}

func (board *Board) InitialState() State {
	var status State

	status.Row = board.StartRow
	status.Col = board.StartCol
	status.CheckpointIdx = 0

	return status
}

// print board dengan posisi player sekarang
func (board *Board) PrintWithActor(curRow, curCol, checkpointIdx int) string {
	hasilTeks := ""

	for i:=0; i < board.N; i++ {
		for j:=0; j < board.M; j++ {
			karakterAsli := board.Tiles[i][j]

			if i == curRow && j == curCol {
				hasilTeks = hasilTeks + "Z"
				continue
			}

			if karakterAsli == 'Z' { // TODO: check if Z or S
				hasilTeks = hasilTeks + "*"
			} else if karakterAsli >= '0' && karakterAsli <= '9' {
				angka := int(karakterAsli - '0')
				diambil := false

				if checkpointIdx > 0 {
					if checkpointIdx >= len(board.Checkpoints) {
						diambil = true
					} else {
						target := board.Checkpoints[checkpointIdx]
						if angka < target {
							diambil = true
						}
					}
				}

				if diambil == true {
					hasilTeks = hasilTeks + "*"
				} else {
					hasilTeks = hasilTeks + string(karakterAsli)
				}

			} else {
				hasilTeks = hasilTeks + string(karakterAsli)
			}
		}

		hasilTeks = hasilTeks + "\n"
	}

	return hasilTeks
}
