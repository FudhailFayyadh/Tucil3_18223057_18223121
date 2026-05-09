package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tucil3/src/puzzle"

	"github.com/eiannone/keyboard"
)

var reader = bufio.NewReader(os.Stdin)

func prompt(msg string) string {
	fmt.Print(msg)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {
	fmt.Println("Ice Sliding Game")

	// cek apakah ada flag --gui waktu program dijalankan
	if len(os.Args) > 1 && os.Args[1] == "--gui" {
		runGUI()
		return
	}

	runCLI()
}

func runCLI() {
	// minta path file input dari user
	filePath := prompt(">> Masukan file input : ")
	board, err := puzzle.ParseBoard(filePath)
	if err != nil {
		fmt.Println("Error membaca file:", err)
		os.Exit(1)
	}

	// pilih algoritma
	algoInput := strings.ToUpper(prompt(">> Algoritma apa yang anda pilih? (UCS/GBFS/A*) : "))

	// kalau GBFS atau A*, tanya heuristik
	var hName string
	if algoInput == "GBFS" || algoInput == "A*" {
		hName = strings.ToUpper(prompt(">> Heuristic apa yang anda pilih? (H1/H2/H3) : "))
	}

	fmt.Println()
	
	// solve
	var result puzzle.SolveResult
	heuristik := puzzle.GetHeuristic(hName)
	var result puzzle.SolveResult

	if algoInput == "UCS" {
		result = puzzle.SolveUCS(board)
	} else if algoInput == "GBFS" {
		result = puzzle.SolveGBFS(board, heuristik)
	} else {
		// def = A* algorithm
		result = puzzle.SolveAStar(board, heuristik)
	}
 
	if !result.Found {
		fmt.Println("Tidak ditemukan solusi.")
		fmt.Printf(">> Waktu eksekusi: %.2f ms\n", result.TimeMs)
		fmt.Printf(">> Banyak iterasi yang dilakukan: %d iterasi\n", result.Iterations)
		return
	}

	// tampilkan rangkuman solusi
	moveStr := movesToString(result.Moves)
	fmt.Printf("Solusi Yang Ditemukan : %s\n", moveStr)
	fmt.Printf("Cost dari Solusi      : %d\n\n", result.TotalCost)

	// tampilkan tiap langkah satu per satu
	fmt.Println("Initial")
	s0 := result.States[0]
	fmt.Print(board.PrintWithActor(s0.Row, s0.Col, s0.CheckpointIdx))

	for i, move := range result.Moves {
		fmt.Printf("\nStep %d : %c\n", i+1, move)
		s := result.States[i+1]
		fmt.Print(board.PrintWithActor(s.Row, s.Col, s.CheckpointIdx))
	}

	fmt.Println()
	fmt.Printf(">> Waktu eksekusi: %.2f ms\n", result.TimeMs)
	fmt.Printf(">> Banyak iterasi yang dilakukan: %d iterasi\n", result.Iterations)

	// replay
	doPlayback := strings.ToLower(prompt("\n>> Apakah Anda ingin melakukan playback? (Ya/Tidak) : "))
	if doPlayback == "ya" || doPlayback == "y" {
		startStepStr := prompt(">> Pada step berapa anda ingin melakukan playback : ")
		startStep, _ := strconv.Atoi(startStepStr)
		if startStep < 0 {
			startStep = 0
		}
		if startStep > len(result.Moves) {
			startStep = len(result.Moves)
		}
		interactivePlayback(board, result, startStep)
	}

	// Save
	doSave := strings.ToLower(prompt("\n>> Apakah Anda ingin menyimpan solusi? (Ya/Tidak) : "))
	if doSave == "ya" || doSave == "y" {
		savePath := prompt(">> Masukan path file output : ")
		if savePath == "" {
			savePath = "solusi.txt"
		}
		err := saveSolution(savePath, board, result, algoInput, hName)
		if err != nil {
			fmt.Println("Error menyimpan:", err)
		} else {
			fmt.Printf(">> Solusi disimpan pada %s\n", savePath)
		}
	}
}

func movesToString(moves []puzzle.Direction) string {
	var sb strings.Builder
	for _, m := range moves {
		sb.WriteByte(byte(m))
	}
	return sb.String()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func interactivePlayback(board *puzzle.Board, result puzzle.SolveResult, startStep int) {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Keyboard tidak tersedia:", err)
		return
	}
	defer keyboard.Close()

	step := startStep
	total := len(result.Moves)

	printStep := func() {
		clearScreen()
		if step == 0 {
			fmt.Println("=== PLAYBACK === Initial State")
		} else {
			fmt.Printf("=== PLAYBACK === Step %d/%d : %c\n", step, total, result.Moves[step-1])
		}
		s := result.States[step]
		fmt.Print(board.PrintWithActor(s.Row, s.Col, s.CheckpointIdx))
		fmt.Printf("\nGunakan ← → untuk navigasi, ESC untuk lompat ke step, Q untuk keluar\n")
	}

	printStep()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}
		switch key {
		case keyboard.KeyArrowRight:
			if step < total {
				step++
			}
			printStep()
		case keyboard.KeyArrowLeft:
			if step > 0 {
				step--
			}
			printStep()
		case keyboard.KeyEsc:
			keyboard.Close()
			fmt.Print("\n>> Lompat ke step : ")
			var input string
			fmt.Scanln(&input)
			n, err := strconv.Atoi(strings.TrimSpace(input))
			if err == nil && n >= 0 && n <= total {
				step = n
			}
			keyboard.Open()
			printStep()
		default:
			if char == 'q' || char == 'Q' {
				return
			}
		}
	}
}

func saveSolution(path string, board *puzzle.Board, result puzzle.SolveResult, algo, heuristic string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "Algoritma : %s", algo)
	if heuristic != "" {
		fmt.Fprintf(w, " dengan heuristik %s", heuristic)
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Solusi    : %s\n", movesToString(result.Moves))
	fmt.Fprintf(w, "Cost      : %d\n", result.TotalCost)
	fmt.Fprintf(w, "Iterasi   : %d\n", result.Iterations)
	fmt.Fprintf(w, "Waktu     : %.2f ms\n\n", result.TimeMs)

	fmt.Fprintln(w, "=== LANGKAH-LANGKAH ===")
	fmt.Fprintln(w, "Initial")
	s0save := result.States[0]
	fmt.Fprint(w, board.PrintWithActor(s0save.Row, s0save.Col, s0save.CheckpointIdx))
	for i, move := range result.Moves {
		s := result.States[i+1]
		fmt.Fprintf(w, "\nStep %d : %c\n", i+1, move)
		fmt.Fprint(w, board.PrintWithActor(s.Row, s.Col, s.CheckpointIdx))
	}

	fmt.Fprintln(w, "\n=== LOG ITERASI ===")
	for _, log := range result.Log {
		fmt.Fprintf(w, "Iter %d: pos=(%d,%d) cpIdx=%d g=%d pri=%d moves=%s\n",
			log.Iteration, log.State.Row, log.State.Col, log.State.CheckpointIdx,
			log.GCost, log.Priority, movesToString(log.Moves))
	}

	return w.Flush()
}
