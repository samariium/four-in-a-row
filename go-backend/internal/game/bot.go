package game

type Bot struct {
	Symbol string // "Y"
}

func (b Bot) opp() string {
	if b.Symbol == "R" { return "Y" }
	return "R"
}

func (b Bot) ChooseMove(g *GameLogic) int {
	// win now
	for c := 0; c < g.Cols; c++ {
		clone := g.Clone()
		if _, ok := clone.DropDisc(c, b.Symbol); ok && clone.CheckWinner(b.Symbol) {
			return c
		}
	}
	// block opp
	opp := b.opp()
	for c := 0; c < g.Cols; c++ {
		clone := g.Clone()
		if _, ok := clone.DropDisc(c, opp); ok && clone.CheckWinner(opp) {
			return c
		}
	}
	// heuristic: center then outwards
	order := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range order {
		if g.ValidColumn(c) { return c }
	}
	for c := 0; c < g.Cols; c++ {
		if g.ValidColumn(c) { return c }
	}
	return 0
}
