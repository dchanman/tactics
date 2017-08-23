package game

import "bytes"

// Square is a (x, y) coordinate tuple
type Square struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Board represents the game board and the pieces on the board
type Board struct {
	Board  []Unit `json:"board"`
	Cols   int    `json:"cols"`
	Rows   int    `json:"rows"`
	nUnits int8
}

// Move is a set of source and destination squares
type Move struct {
	Src Square
	Dst Square
}

// Step assumes Src and Dst are adjacent squares
type Step Move

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

func newBoard(cols int, rows int) Board {
	b := Board{
		Board: make([]Unit, cols*rows),
		Cols:  cols,
		Rows:  rows}
	return b
}

func (b *Board) isValid(x int, y int) bool {
	return (x >= 0 && y >= 0 && x < b.Cols && y < b.Rows)
}

func (b *Board) get(x int, y int) Unit {
	if !b.isValid(x, y) {
		return Unit{Exists: false}
	}
	return b.Board[x*b.Rows+y]
}

func (b *Board) set(x int, y int, u Unit) {
	if !b.isValid(x, y) {
		return
	}
	prev := b.get(x, y)
	if prev.Exists {
		b.nUnits -= prev.Stack
	}
	b.Board[x*b.Rows+y] = u
	if u.Exists {
		b.nUnits += u.Stack
	}
}

func (b *Board) pickup(x int, y int) Unit {
	ret := b.get(x, y)
	b.set(x, y, Unit{Exists: false})
	return ret
}

func (b *Board) getLineInDirection(dir func(s *Square) Square, origin Square) []Square {
	ret := make([]Square, 0)
	for next := dir(&origin); b.isValid(next.X, next.Y); next = dir(&next) {
		ret = append(ret, next)
	}
	return ret
}

func (b *Board) getValidMoves(x int, y int) []Square {
	u := b.get(x, y)
	if u.Exists == false {
		return make([]Square, 0)
	}
	return u.getValidMoves(b, Square{x, y})
}

func computeIncrement(src int, dst int) int {
	if src > dst {
		return -1
	}
	return 1
}

func decomposeMoveToSteps(m Move) []Step {
	ret := make([]Step, 0)
	if m.Src.X == m.Dst.X {
		inc := computeIncrement(m.Src.Y, m.Dst.Y)
		last := m.Src
		for i := m.Src.Y + inc; last != m.Dst; i += inc {
			next := Square{X: m.Src.X, Y: i}
			ret = append(ret, Step{Src: last, Dst: next})
			last = next
		}
	} else if m.Src.Y == m.Dst.Y {
		inc := computeIncrement(m.Src.X, m.Dst.X)
		last := m.Src
		for i := m.Src.X + inc; last != m.Dst; i += inc {
			next := Square{X: i, Y: m.Src.Y}
			ret = append(ret, Step{Src: last, Dst: next})
			last = next
		}
	}
	return ret
}

func (b *Board) resolveStep(step1 Step, step2 Step) (bool, bool) {
	u1 := b.pickup(step1.Src.X, step1.Src.Y)
	u2 := b.pickup(step2.Src.X, step2.Src.Y)
	// Case 1: Both steps converge onto the same square
	if step1.Dst == step2.Dst {
		udst := b.pickup(step1.Dst.X, step1.Dst.Y)
		b.set(step1.Dst.X, step1.Dst.Y, stack(u1, udst))
		udst = b.pickup(step2.Dst.X, step2.Dst.Y)
		b.set(step2.Dst.X, step2.Dst.Y, stack(u2, udst))
		return true, true
	}
	// Case 2: Steps are on adjacent squares, moving into one another
	// In this case, the stacking is applied, but neither piece will move
	if step1.Dst == step2.Src && step1.Src == step2.Dst {
		collision := stack(u1, u2)
		if collision.Exists && collision.Team == u1.Team {
			b.set(step1.Src.X, step1.Src.Y, collision)
		} else if collision.Exists && collision.Team == u2.Team {
			b.set(step2.Src.X, step2.Src.Y, collision)
		}
		return true, true
	}
	// Case 3: Two exclusive moves
	collision1 := b.get(step1.Dst.X, step1.Dst.Y).Exists
	collision2 := b.get(step2.Dst.X, step2.Dst.Y).Exists
	b.set(step1.Dst.X, step1.Dst.Y, stack(u1, b.get(step1.Dst.X, step1.Dst.Y)))
	b.set(step2.Dst.X, step2.Dst.Y, stack(u2, b.get(step2.Dst.X, step2.Dst.Y)))
	return collision1, collision2
}

func checkWinCondition(b *Board) (bool, Team) {
	if b.nUnits <= 0 {
		return true, 0
	}
	team1win := int8(0)
	team2win := int8(0)
	// Let rank 0 be team 1's "endzone"
	for i := 0; i < b.Cols; i++ {
		if b.get(i, 0).Exists && b.get(i, 0).Team == 1 {
			team1win = b.get(i, 0).Stack
		}
	}
	// Let rank nRows be team 2's "endzone"
	for i := 0; i < b.Cols; i++ {
		if b.get(i, b.Rows-1).Exists && b.get(i, b.Rows-1).Team == 2 {
			team2win = b.get(i, b.Rows-1).Stack
		}
	}
	if team1win > team2win {
		return true, 1
	} else if team2win > team1win {
		return true, 2
	} else if team1win > 0 {
		// draw game: both sides got the same number of pieces across
		return true, 0
	}
	return false, 0
}

func (b *Board) resolveMove(move1 Move, move2 Move) (bool, Team) {
	// TODO: validate moves
	// logrus.WithFields(logrus.Fields{"Board": b}).Info("init")
	// defer logrus.WithFields(logrus.Fields{"Board": b}).Info("fini")
	steps1 := decomposeMoveToSteps(move1)
	steps2 := decomposeMoveToSteps(move2)

	nSteps := len(steps1)
	if len(steps2) > len(steps1) {
		nSteps = len(steps2)
	}

	var s1 Step
	var s2 Step
	stopped1 := false
	stopped2 := false
	for i := 0; i < nSteps && !(stopped1 && stopped2); i++ {
		if stopped1 || i >= len(steps1) {
			s1 = Step{Src: move1.Dst, Dst: move1.Dst}
		} else {
			s1 = steps1[i]
		}
		if stopped2 || i >= len(steps2) {
			s2 = Step{Src: move2.Dst, Dst: move2.Dst}
		} else {
			s2 = steps2[i]
		}
		collision1, collision2 := b.resolveStep(s1, s2)
		stopped1 = collision1 || stopped1
		stopped2 = collision2 || stopped2
		winner, team := checkWinCondition(b)
		if winner {
			return true, team
		}
	}
	return false, 0
}

func (b *Board) String() string {
	var buf bytes.Buffer

	for i := 0; i < len(b.Board); i++ {
		if i%b.Rows == 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(b.Board[i].String())
		buf.WriteString("\t\t| ")
	}
	return buf.String()
}
