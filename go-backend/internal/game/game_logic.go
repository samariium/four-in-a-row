package game

type GameLogic struct {
	Rows  int
	Cols  int
	Board [][]*string // nil or "R"/"Y"
}

func NewGame() *GameLogic {
	rows, cols := 6, 7
	b := make([][]*string, rows)
	for r := 0; r < rows; r++ {
		b[r] = make([]*string, cols)
	}
	return &GameLogic{Rows: rows, Cols: cols, Board: b}
}

func (g *GameLogic) Clone() *GameLogic {
	n := NewGame()
	for r := range g.Board {
		for c := range g.Board[r] {
			if g.Board[r][c] != nil {
				val := *g.Board[r][c]
				n.Board[r][c] = &val
			}
		}
	}
	return n
}

func (g *GameLogic) ValidColumn(col int) bool {
	return col >= 0 && col < g.Cols && g.Board[0][col] == nil
}

func (g *GameLogic) DropDisc(col int, player string) (row int, ok bool) {
	if !g.ValidColumn(col) { return -1, false }
	for r := g.Rows - 1; r >= 0; r-- {
		if g.Board[r][col] == nil {
			val := player
			g.Board[r][col] = &val
			return r, true
		}
	}
	return -1, false
}

func (g *GameLogic) CheckWinner(p string) bool {
	b := g.Board; R, C := g.Rows, g.Cols
	// horizontal
	for r := 0; r < R; r++ {
		for c := 0; c <= C-4; c++ {
			if b[r][c] != nil && b[r][c+1] != nil && b[r][c+2] != nil && b[r][c+3] != nil &&
				*b[r][c] == p && *b[r][c+1] == p && *b[r][c+2] == p && *b[r][c+3] == p {
				return true
			}
		}
	}
	// vertical
	for c := 0; c < C; c++ {
		for r := 0; r <= R-4; r++ {
			if b[r][c] != nil && b[r+1][c] != nil && b[r+2][c] != nil && b[r+3][c] != nil &&
				*b[r][c] == p && *b[r+1][c] == p && *b[r+2][c] == p && *b[r+3][c] == p {
				return true
			}
		}
	}
	// diag down-right
	for r := 0; r <= R-4; r++ {
		for c := 0; c <= C-4; c++ {
			if b[r][c] != nil && b[r+1][c+1] != nil && b[r+2][c+2] != nil && b[r+3][c+3] != nil &&
				*b[r][c] == p && *b[r+1][c+1] == p && *b[r+2][c+2] == p && *b[r+3][c+3] == p {
				return true
			}
		}
	}
	// diag up-right
	for r := 3; r < R; r++ {
		for c := 0; c <= C-4; c++ {
			if b[r][c] != nil && b[r-1][c+1] != nil && b[r-2][c+2] != nil && b[r-3][c+3] != nil &&
				*b[r][c] == p && *b[r-1][c+1] == p && *b[r-2][c+2] == p && *b[r-3][c+3] == p {
				return true
			}
		}
	}
	return false
}

func (g *GameLogic) IsFull() bool {
	for c := 0; c < g.Cols; c++ {
		if g.Board[0][c] == nil { return false }
	}
	return true
}
