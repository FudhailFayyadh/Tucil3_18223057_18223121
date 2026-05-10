# Ice Sliding Puzzle Solver

Implementasi solver untuk permainan Ice Sliding Puzzle menggunakan algoritma pathfinding UCS, GBFS, dan A* dalam bahasa Go.

## Deskripsi

Program ini menyelesaikan Ice Sliding Puzzle di mana pemain harus menggerakkan karakter dari titik awal menuju titik tujuan di atas permukaan es yang licin. Karakter tidak berhenti bergerak sampai menabrak dinding atau rintangan.

**Algoritma yang diimplementasikan:**
- UCS (Uniform Cost Search)
- GBFS (Greedy Best First Search)
- A* Search

**Heuristik (untuk GBFS dan A*):**
- H1: Manhattan Distance (|Δrow| + |Δcol|)
- H2: Euclidean Distance (√(Δrow² + Δcol²))
- H3: Chebyshev Distance (max(|Δrow|, |Δcol|))

**Bonus:**
- GUI dengan Fyne
- 2 heuristik tambahan (H2 dan H3)

## Requirement

### Windows
- Go 1.21+
- GCC (via MSYS2 UCRT64): `pacman -S mingw-w64-ucrt-x86_64-gcc`
- MSYS2 ucrt64/bin harus ada di PATH

### Linux
```bash
sudo apt install gcc libgl1-mesa-dev xorg-dev
```

## Cara Kompilasi

```bash
go mod tidy
go build -o bin/tucil3.exe ./src/    # Windows
go build -o bin/tucil3 ./src/        # Linux
```

## Cara Menjalankan

### Mode CLI
```bash
./bin/tucil3.exe
# atau
./bin/tucil3
```

Ikuti prompt:
1. Masukkan path file input (.txt)
2. Pilih algoritma (UCS/GBFS/A*)
3. Pilih heuristik jika memilih GBFS atau A* (H1/H2/H3)
4. Program menampilkan solusi dan visualisasi langkah demi langkah
5. Opsi playback interaktif dengan tombol ← →
6. Opsi menyimpan solusi ke file .txt

### Mode GUI
```bash
./bin/tucil3.exe --gui
```

## Format File Input

```
N M
[N baris tile]
[N baris cost]
```

**Keterangan tile:**
- `*` = path yang bisa dilewati
- `X` = rintangan/batu
- `L` = lava (game over jika dilewati)
- `Z` = posisi aktor (start)
- `O` = titik tujuan
- `0-9` = checkpoint yang harus dilewati sesuai urutan

## Contoh Input

```
7 7
XXXXXXX
X0****X
X**X**X
X****OX
X***1LX
XZ**X*X
XXXXXXX
999 999 999 999 999 999 999
999 3 5 2 8 1 999
999 7 4 999 6 9 999
999 2 8 3 5 4 999
999 6 1 7 2 999 999
999 9 3 4 999 8 999
999 999 999 999 999 999 999
```

## Struktur Repository

```
Tucil3_18223057_18223121/
├── src/
│   ├── main.go          # Entry point CLI
│   ├── gui.go           # GUI (Fyne)
│   └── puzzle/
│       ├── board.go     # Parsing board
│       ├── state.go     # State + sliding simulation
│       ├── heuristic.go # H1, H2, H3
│       └── solver.go    # UCS, GBFS, A*
├── test/                # Test cases
└── README.md
```

## Author

Fudhail Fayyadh 18223121
Stanislaus Ardy Bramantyo 18223057