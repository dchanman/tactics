package main

import (
	"game"

	"github.com/gopherjs/gopherjs/js"
)

func newBoard(cols int, rows int) *js.Object {
	b := game.NewBoard(cols, rows)
	return js.MakeWrapper(&b)
}

func unit(team game.Team, stack int) *js.Object {
	return js.MakeWrapper(game.Unit{Team: team, Stack: stack, Exists: true})
}

func newMove(x0 int, y0 int, x1 int, y1 int) *js.Object {
	return js.MakeWrapper(game.Move{game.Square{x0, y0}, game.Square{x1, y1}})
}

func main() {
	js.Global.Set("Engine", map[string]interface{}{
		"NewBoard": newBoard,
		"NewUnit":  unit,
		"NewMove":  newMove})
}
