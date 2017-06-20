package main

// Game is the main game engine
type Game struct {
	b board
}

func NewGame() Game {
	return Game{b: NewBoard(10, 10)}
}

func (g *Game) GetStateJSON() string {
	return g.b.ToJSON()
}
