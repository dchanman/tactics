package game

// unit is a basic unit in the game
type Unit struct {
	Name   string `json:"name,omitempty"`
	Class  string `json:"class,omitempty"`
	Team   int8   `json:"team,omitempty"`
	Stack  int8   `json:"stack,omitempty"`
	Exists bool   `json:"exists"`
}
