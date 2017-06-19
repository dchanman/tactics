package game

import "github.com/sirupsen/logrus"

type square struct {
	x int
	y int
}

type board struct {
	board []unit
	cols  int
	rows  int
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

func newBoard(cols int, rows int) board {
	b := board{
		board: make([]unit, cols*rows),
		cols:  cols,
		rows:  rows}
	return b
}

func (b *board) isValid(x int, y int) bool {
	return (x >= 0 && y >= 0 && x < b.cols && y < b.rows)
}

func (b *board) get(x int, y int) unit {
	if !b.isValid(x, y) {
		log.WithFields(logrus.Fields{
			"x":    x,
			"y":    y,
			"cols": b.cols,
			"rows": b.rows}).Error("Index out of bounds")
		return unit{exists: false}
	}
	return b.board[x*b.cols+y]
}

func (b *board) set(x int, y int, u unit) {
	if !b.isValid(x, y) {
		log.WithFields(logrus.Fields{
			"x":    x,
			"y":    y,
			"cols": b.cols,
			"rows": b.rows}).Error("Index out of bounds")
		return
	}
	b.board[x*b.cols+y] = u
}

func (b *board) getValidMoves(x int, y int) []square {
	u := b.get(x, y)
	if u.mov == 0 {
		return make([]square, 0)
	}
	queue := make([][]square, 0)
	for i := 0; i <= int(u.mov); i++ {
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
			moves = append(moves, curr)
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
	return moves[1:]
}
