package main

import (
	"game"

	"github.com/gopherjs/gopherjs/js"
)

func newBoardFromBoard(cols int, rows int, obj []struct {
	*js.Object
	Team   int  `js:"team"`
	Stack  int  `js:"stack"`
	Exists bool `js:"exists"`
}) *js.Object {
	b := game.NewBoard(cols, rows)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			i := (x * rows) + y
			u := game.Unit{Team: game.Team(obj[i].Team), Stack: obj[i].Stack, Exists: obj[i].Exists}
			b.Set(x, y, u)
		}
	}
	return js.MakeWrapper(&b)
}

func newBoard(cols int, rows int) *js.Object {
	b := game.NewBoard(cols, rows)
	return js.MakeWrapper(&b)
}

func unit(team game.Team, stack int) *js.Object {
	return js.MakeWrapper(game.Unit{Team: team, Stack: stack, Exists: true})
}

func newMove(x0 int, y0 int, x1 int, y1 int) *js.Object {
	return js.MakeWrapper(game.Move{Src: game.Square{X: x0, Y: y0}, Dst: game.Square{X: x1, Y: y1}})
}

func main() {
	js.Global.Set("Engine", map[string]interface{}{
		"NewBoard":          newBoard,
		"NewBoardFromBoard": newBoardFromBoard,
		"NewUnit":           unit,
		"NewMove":           newMove})
}
