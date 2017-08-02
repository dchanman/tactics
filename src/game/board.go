package game

import "encoding/json"

type Square struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Board struct {
	Board []Unit `json:"board"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
}

type moveSearch struct {
	Square
	movRemainder int8
}

func (s *Square) up() Square {
	return Square{s.X, s.Y - 1}
}
func (s *Square) down() Square {
	return Square{s.X, s.Y + 1}
}
func (s *Square) left() Square {
	return Square{s.X - 1, s.Y}
}
func (s *Square) right() Square {
	return Square{s.X + 1, s.Y}
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
	return b.Board[x*b.Rows+y]
}

func (b *Board) Set(x int, y int, u Unit) {
	if !b.isValid(x, y) {
		return
	}
	b.Board[x*b.Rows+y] = u
}

func (b *Board) GetValidMoves(x int, y int) []Square {
	u := b.Get(x, y)
	if u.Exists == false {
		return make([]Square, 0)
	}
	moves := make([]Square, 0)
	searchDirHelper := func(dir func(s *Square) Square, origin Square) []Square {
		ret := make([]Square, 0)
		for next := dir(&origin); b.isValid(next.X, next.Y); next = dir(&next) {
			ret = append(ret, next)
		}
		return ret
	}
	moves = append(moves, searchDirHelper((*Square).up, Square{x, y})...)
	moves = append(moves, searchDirHelper((*Square).down, Square{x, y})...)
	moves = append(moves, searchDirHelper((*Square).left, Square{x, y})...)
	moves = append(moves, searchDirHelper((*Square).right, Square{x, y})...)
	return moves
}

func (b *Board) ToJSON() string {
	marshalled, _ := json.Marshal(b)
	log.Info(string(marshalled))
	return string(marshalled)
}
