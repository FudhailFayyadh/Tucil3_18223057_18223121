package puzzle

import "math"

// tipe fungsi heuristik buat A*
type HeuristicFn func(s State, b *Board) int

// heuristik 1: manhattan
func H1Manhattan(s State, b *Board) int {
	selisihBaris := s.Row - b.GoalRow
	if selisihBaris < 0 {
		selisihBaris = -selisihBaris 
	}

	selisihKolom := s.Col - b.GoalCol
	if selisihKolom < 0 {
		selisihKolom = -selisihBaris
	}

	hasilAkhir := selisihBaris + selisihKolom
	return hasilAkhir
}

// heuristik 2: euclidean
func H2Euclidean(s State, b *Board) int {
	selisihBaris := s.Row - b.GoalRow
	selisihKolom := s.Col - b.GoalCol

	floatSB := float64(selisihBaris)
	floatSK := float64(selisihKolom)

	kuadratSB := floatSB * floatSB
	kuadratSK := floatSK * floatSK

	// membulatkan ke bawah
	total := kuadratSB+kuadratSK
	rootTotal := math.Sqrt(total)
	hasilAkhir := int (rootTotal)
	
	return hasilAkhir
}

// heuristic 3: chebyshev
func H3Chebyshev(s State, b *Board) int {
	selisihBaris := s.Row - b.GoalRow
	if selisihBaris < 0 {
		selisihBaris = -selisihBaris 
	}

	selisihKolom := s.Col - b.GoalCol
	if selisihKolom < 0 {
		selisihKolom = -selisihBaris
	}

	var hasilAkhir int

	if selisihBaris > selisihKolom {
		hasilAkhir = selisihBaris
	}
	else {
		hasilAkhir = selisihKolom
	}

	return hasilAkhir
}

func GetHeuristic(name string) HeuristicFn {
	var hasil HeuristicFn

	if name == "H1" {
		hasil = H1Manhattan
	}
	else if name == "H2" {
		hasil = H2Euclidean
	}
	else if name == "H3" {
		hasil = H3Chebyshev
	}
	else {
		hasil = H1Manhattan // Fail Haven
	}

	return hasil
}
