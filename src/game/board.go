package game

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type square struct {
	x int
	y int
}

type Board struct {
	Board []Unit
	Cols  int
	Rows  int
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
		log.WithFields(logrus.Fields{
			"x":    x,
			"y":    y,
			"Cols": b.Cols,
			"Rows": b.Rows}).Error("Index out of bounds")
		return Unit{Exists: false}
	}
	return b.Board[x*b.Cols+y]
}

func (b *Board) Set(x int, y int, u Unit) {
	if !b.isValid(x, y) {
		log.WithFields(logrus.Fields{
			"x":    x,
			"y":    y,
			"Cols": b.Cols,
			"Rows": b.Rows}).Error("Index out of bounds")
		return
	}
	b.Board[x*b.Cols+y] = u
}

func (b *Board) getValidMoves(x int, y int) []square {
	u := b.Get(x, y)
	if u.Mov == 0 {
		return make([]square, 0)
	}
	queue := make([][]square, 0)
	for i := 0; i <= int(u.Mov); i++ {
		queue = append(queue, make([]square, 0))
	}
	queue[0] = append(queue[0], square{x, y})
	visited := make(map[square]bool)
	moves := make([]square, 0)
	for i := 0; i < len(queue); i++ {
		for j := 0; j < len(queue[i]); j++ {
			curr := queue[i][j]
			if visited[curr] || !b.isValid(curr.x, curr.y) {
				continue
			}
			visited[curr] = true
			o := b.Get(curr.x, curr.y)
			if !o.Exists {
				moves = append(moves, curr)
			} else if o.Team != u.Team {
				// Enemy pieces "block" movement
				continue
			}
			if i+1 < len(queue) {
				if !visited[curr.left()] {
					queue[i+1] = append(queue[i+1], curr.left())
				}
				if !visited[curr.right()] {
					queue[i+1] = append(queue[i+1], curr.right())
				}
				if !visited[curr.up()] {
					queue[i+1] = append(queue[i+1], curr.up())
				}
				if !visited[curr.down()] {
					queue[i+1] = append(queue[i+1], curr.down())
				}
			}
		}
	}
	return moves
}

func (b *Board) ToJSON() string {
	marshalled, _ := json.Marshal(b)
	log.Info(string(marshalled))
	return string(marshalled)
}
