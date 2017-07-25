package game

// Game is the main game engine
type Game struct {
	B Board `json:"board,omitempty"`
}

func NewGame() Game {
	return Game{B: NewBoard(10, 10)}
}
