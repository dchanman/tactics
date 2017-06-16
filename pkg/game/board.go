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
	moves := make([]square, 0)
	visited := make(map[square]bool)
	tovisit := make([]moveSearch, 0)
	i := 0
	u := b.get(x, y)
	tovisit = append(tovisit, moveSearch{square{x, y}, u.mov})

	for i < len(tovisit) {
		sq := tovisit[i].square
		if !visited[sq] && b.isValid(sq.x, sq.y) && tovisit[i].movRemainder >= 0 {
			visited[tovisit[i].square] = true
			rem := tovisit[i].movRemainder - 1
			moves = append(moves, tovisit[i].square)
			// Add right move
			tovisit = append(tovisit, moveSearch{square{sq.x + 1, sq.y}, rem})
			// Add left move
			tovisit = append(tovisit, moveSearch{square{sq.x - 1, sq.y}, rem})
			// Add up move
			tovisit = append(tovisit, moveSearch{square{sq.x, sq.y - 1}, rem})
			// Add down move
			tovisit = append(tovisit, moveSearch{square{sq.x, sq.y + 1}, rem})
		}
		i++
	}

	return moves[1:]
}
