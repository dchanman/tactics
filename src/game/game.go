package game

// Game is the main game engine
type Game struct {
	B Board
}

func NewGame() Game {
	return Game{B: NewBoard(10, 10)}
}
