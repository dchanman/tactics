package game

import "encoding/json"

type square struct {
	x int
	y int
}

type Board struct {
	Board []Unit `json:"board"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
}

type moveSearch struct {
	square
	movRemainder int8
}

func (s *square) up() square {
	return square{s.x, s.y - 1}
}
func (s *square) down() square {
	return square{s.x, s.y + 1}
}
func (s *square) left() square {
	return square{s.x - 1, s.y}
}
func (s *square) right() square {
	return square{s.x + 1, s.y}
}

func NewBoard(cols int, rows int) Board {
	b := Board{
		Board: make([]Unit, cols*rows),
		Cols:  cols,
		Rows:  rows}
	return b
}

func (b *Board) isValid(x int, y int) bool {
	return (x >= 0 && y >= 0 && x < b.Cols && y < b.Rows)
}

func (b *Board) Get(x int, y int) Unit {
	if !b.isValid(x, y) {
		return Unit{Exists: false}
	}
	return b.Board[x*b.Cols+y]
}

func (b *Board) Set(x int, y int, u Unit) {
	if !b.isValid(x, y) {
		return
	}
	b.Board[x*b.Cols+y] = u
}

func (b *Board) getValidMoves(x int, y int) []square {
	u := b.Get(x, y)
	if u.Exists == false {
		return make([]square, 0)
	}
	moves := make([]square, 0)
	searchDirHelper := func(dir func(s *square) square, origin square) []square {
		ret := make([]square, 0)
		for next := dir(&origin); b.isValid(next.x, next.y); next = dir(&next) {
			ret = append(ret, next)
		}
		return ret
	}
	moves = append(moves, searchDirHelper((*square).up, square{x, y})...)
	moves = append(moves, searchDirHelper((*square).down, square{x, y})...)
	moves = append(moves, searchDirHelper((*square).left, square{x, y})...)
	moves = append(moves, searchDirHelper((*square).right, square{x, y})...)
	return moves
}

func (b *Board) ToJSON() string {
	marshalled, _ := json.Marshal(b)
	log.Info(string(marshalled))
	return string(marshalled)
}
