package game

import "github.com/sirupsen/logrus"

type board struct {
	board []unit
	cols  int
	rows  int
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
