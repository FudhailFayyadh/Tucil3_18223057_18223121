package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"
	"tucil3/src/puzzle"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	colorWall       = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	colorPath       = color.RGBA{R: 60, G: 60, B: 70, A: 255}
	colorLava       = color.RGBA{R: 200, G: 50, B: 50, A: 255}
	colorGoal       = color.RGBA{R: 50, G: 180, B: 80, A: 255}
	colorActor      = color.RGBA{R: 50, G: 100, B: 220, A: 255}
	colorCheckpoint = color.RGBA{R: 220, G: 180, B: 40, A: 255}
	colorVisited    = color.RGBA{R: 100, G: 140, B: 200, A: 180}
)

type guiState struct {
	board       *puzzle.Board
	result      *puzzle.SolveResult
	currentStep int
	totalSteps  int
	ticker      *time.Ticker
	stopTicker  chan struct{}
}

func runGUI() {
	a := app.New()
	w := a.NewWindow("Ice Sliding Puzzle Solver")
	w.Resize(fyne.NewSize(900, 650))

	gs := &guiState{}

	// --- Controls panel ---
	fileLabel := widget.NewLabel("Belum ada file dipilih")
	fileLabel.Wrapping = fyne.TextWrapWord

	algoSelect := widget.NewSelect([]string{"UCS", "GBFS", "A*"}, nil)
	algoSelect.SetSelected("A*")

	heurLabel := widget.NewLabel("Heuristik:")
	heurSelect := widget.NewSelect([]string{"H1 (Manhattan)", "H2 (Euclidean)", "H3 (Chebyshev)"}, nil)
	heurSelect.SetSelected("H1 (Manhattan)")
	heurLabel.Hide()
	heurSelect.Hide()

	algoSelect.OnChanged = func(s string) {
		if s == "GBFS" || s == "A*" {
			heurLabel.Show()
			heurSelect.Show()
		} else {
			heurLabel.Hide()
			heurSelect.Hide()
		}
	}

	resultLabel := widget.NewLabel("")
	resultLabel.Wrapping = fyne.TextWrapWord

	// --- Board canvas ---
	boardContainer := container.NewWithoutLayout()
	boardScroll := container.NewScroll(boardContainer)
	boardScroll.SetMinSize(fyne.NewSize(500, 400))

	// --- Playback controls ---
	stepLabel := widget.NewLabel("Step 0/0")
	prevBtn := widget.NewButton("◀", nil)
	nextBtn := widget.NewButton("▶", nil)
	speedSlider := widget.NewSlider(0.1, 3.0)
	speedSlider.Value = 1.0
	speedSlider.Step = 0.1
	autoPlayBtn := widget.NewButton("▶ Auto", nil)
	saveBtn := widget.NewButton("Simpan Solusi", nil)
	saveBtn.Disable()

	prevBtn.Disable()
	nextBtn.Disable()

	var drawBoard func()

	drawBoard = func() {
		if gs.board == nil {
			return
		}
		boardContainer.Objects = nil
		cellSize := float32(40)
		b := gs.board

		// visited set for current solution up to currentStep
		visitedCells := map[[2]int]bool{}
		if gs.result != nil {
			for i := 1; i <= gs.currentStep && i < len(gs.result.States); i++ {
				s := gs.result.States[i]
				visitedCells[[2]int{s.Row, s.Col}] = true
			}
		}

		var actorRow, actorCol, cpIdx int
		if gs.result != nil && gs.currentStep < len(gs.result.States) {
			s := gs.result.States[gs.currentStep]
			actorRow, actorCol, cpIdx = s.Row, s.Col, s.CheckpointIdx
		} else {
			actorRow, actorCol = b.StartRow, b.StartCol
		}

		for i := 0; i < b.N; i++ {
			for j := 0; j < b.M; j++ {
				x := float32(j) * cellSize
				y := float32(i) * cellSize

				tile := b.Tiles[i][j]
				var bg color.Color
				var label string

				switch tile {
				case puzzle.TileWall:
					bg = colorWall
					label = "X"
				case puzzle.TileLava:
					bg = colorLava
					label = "L"
				case puzzle.TileGoal:
					bg = colorGoal
					label = "O"
				case puzzle.TileStart:
					bg = colorPath
					label = "*"
				default:
					if tile >= '0' && tile <= '9' {
						digit := int(tile - '0')
						collected := cpIdx > 0 && (cpIdx >= len(b.Checkpoints) || digit < b.Checkpoints[cpIdx])
						if collected {
							bg = colorPath
							label = "*"
						} else {
							bg = colorCheckpoint
							label = string(tile)
						}
					} else {
						bg = colorPath
						label = " "
					}
				}

				if i == actorRow && j == actorCol {
					bg = colorActor
					label = "Z"
				} else if visitedCells[[2]int{i, j}] && tile != puzzle.TileWall {
					bg = colorVisited
				}

				rect := canvas.NewRectangle(bg)
				rect.SetMinSize(fyne.NewSize(cellSize, cellSize))
				rect.Move(fyne.NewPos(x, y))
				rect.Resize(fyne.NewSize(cellSize-1, cellSize-1))

				txt := canvas.NewText(label, color.White)
				txt.TextSize = 14
				txt.Alignment = fyne.TextAlignCenter
				txt.Move(fyne.NewPos(x+cellSize/4, y+cellSize/4))

				boardContainer.Add(rect)
				boardContainer.Add(txt)
			}
		}
		boardContainer.Resize(fyne.NewSize(
			float32(b.M)*cellSize,
			float32(b.N)*cellSize,
		))
		boardContainer.Refresh()
	}

	updateStep := func() {
		if gs.result == nil {
			return
		}
		total := gs.totalSteps
		cur := gs.currentStep
		if cur == 0 {
			stepLabel.SetText(fmt.Sprintf("Step 0/%d (Initial)", total))
		} else {
			stepLabel.SetText(fmt.Sprintf("Step %d/%d (%c)", cur, total, gs.result.Moves[cur-1]))
		}
		drawBoard()
	}

	prevBtn.OnTapped = func() {
		if gs.currentStep > 0 {
			gs.currentStep--
			updateStep()
		}
	}
	nextBtn.OnTapped = func() {
		if gs.currentStep < gs.totalSteps {
			gs.currentStep++
			updateStep()
		}
	}

	var stopAutoPlay func()
	stopAutoPlay = func() {
		if gs.ticker != nil {
			gs.ticker.Stop()
			gs.ticker = nil
		}
		if gs.stopTicker != nil {
			select {
			case gs.stopTicker <- struct{}{}:
			default:
			}
		}
		autoPlayBtn.SetText("▶ Auto")
	}

	autoPlayBtn.OnTapped = func() {
		if gs.ticker != nil {
			stopAutoPlay()
			return
		}
		if gs.result == nil {
			return
		}
		interval := time.Duration(float64(time.Second) / speedSlider.Value)
		gs.ticker = time.NewTicker(interval)
		gs.stopTicker = make(chan struct{}, 1)
		autoPlayBtn.SetText("⏹ Stop")
		go func() {
			for {
				select {
				case <-gs.ticker.C:
					if gs.currentStep < gs.totalSteps {
						gs.currentStep++
						updateStep()
					} else {
						stopAutoPlay()
						return
					}
				case <-gs.stopTicker:
					return
				}
			}
		}()
	}

	saveBtn.OnTapped = func() {
		if gs.result == nil || gs.board == nil {
			return
		}
		dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			path := uc.URI().Path()
			algo := algoSelect.Selected
			hName := heurName(heurSelect.Selected)
			saveSolution(path, gs.board, *gs.result, algo, hName)
		}, w)
	}

	// Run button
	runBtn := widget.NewButton("▶ Jalankan Solver", func() {
		if gs.board == nil {
			dialog.ShowInformation("Error", "Pilih file terlebih dahulu", w)
			return
		}
		stopAutoPlay()

		algo := algoSelect.Selected
		hName := heurName(heurSelect.Selected)
		h := puzzle.GetHeuristic(hName)

		resultLabel.SetText("Mencari solusi...")

		go func() {
			var res puzzle.SolveResult
			switch algo {
			case "UCS":
				res = puzzle.SolveUCS(gs.board)
			case "GBFS":
				res = puzzle.SolveGBFS(gs.board, h)
			default:
				res = puzzle.SolveAStar(gs.board, h)
			}

			gs.result = &res
			gs.currentStep = 0
			gs.totalSteps = len(res.Moves)

			if res.Found {
				moveStr := movesToString(res.Moves)
				info := fmt.Sprintf("Solusi: %s\nCost: %d\nIterasi: %d\nWaktu: %.2f ms",
					moveStr, res.TotalCost, res.Iterations, res.TimeMs)
				resultLabel.SetText(info)
				prevBtn.Enable()
				nextBtn.Enable()
				saveBtn.Enable()
			} else {
				resultLabel.SetText(fmt.Sprintf("Tidak ditemukan solusi.\nIterasi: %d\nWaktu: %.2f ms",
					res.Iterations, res.TimeMs))
			}
			updateStep()
		}()
	})

	// Load file button
	loadBtn := widget.NewButton("Pilih File", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			path := uc.URI().Path()
			uc.Close()

			// On Windows the path may have a leading slash
			if len(path) > 2 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}

			b, err := puzzle.ParseBoard(path)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			gs.board = b
			gs.result = nil
			gs.currentStep = 0
			gs.totalSteps = 0
			fileLabel.SetText(fmt.Sprintf("File: %s (%dx%d)", path, b.N, b.M))
			resultLabel.SetText("")
			stepLabel.SetText("Step 0/0")
			prevBtn.Disable()
			nextBtn.Disable()
			saveBtn.Disable()
			drawBoard()
		}, w)
	})

	// Layout
	controls := container.NewVBox(
		widget.NewLabel("=== Ice Sliding Puzzle Solver ==="),
		loadBtn,
		fileLabel,
		widget.NewSeparator(),
		widget.NewLabel("Algoritma:"),
		algoSelect,
		heurLabel,
		heurSelect,
		runBtn,
		widget.NewSeparator(),
		resultLabel,
		widget.NewSeparator(),
		saveBtn,
	)
	controlsScroll := container.NewScroll(controls)
	controlsScroll.SetMinSize(fyne.NewSize(220, 600))

	playbackBar := container.NewHBox(
		prevBtn,
		stepLabel,
		nextBtn,
		layout.NewSpacer(),
		widget.NewLabel("Speed:"),
		speedSlider,
		autoPlayBtn,
	)

	rightPanel := container.NewBorder(nil, playbackBar, nil, nil, boardScroll)
	split := container.NewHSplit(controlsScroll, rightPanel)
	split.SetOffset(0.25)

	w.SetContent(split)
	w.ShowAndRun()
}

func heurName(selected string) string {
	switch {
	case strings.HasPrefix(selected, "H2"):
		return "H2"
	case strings.HasPrefix(selected, "H3"):
		return "H3"
	default:
		return "H1"
	}
}

