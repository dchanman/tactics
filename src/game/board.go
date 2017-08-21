package game

import "bytes"

type Square struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Board struct {
	Board []Unit `json:"board"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
}

type Move struct {
	Src Square
	Dst Square
}

// Assumption: A step's Src and Dst are adjacent squares
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

func (b *Board) pickup(x int, y int) Unit {
	ret := b.Get(x, y)
	b.Set(x, y, Unit{Exists: false})
	return ret
}

func (b *Board) getValidMovesHelper(dir func(s *Square) Square, origin Square) []Square {
	ret := make([]Square, 0)
	for next := dir(&origin); b.isValid(next.X, next.Y); next = dir(&next) {
		ret = append(ret, next)
	}
	return ret
}

func (b *Board) getValidMoves(x int, y int) []Square {
	u := b.Get(x, y)
	if u.Exists == false {
		return make([]Square, 0)
	}
	moves := make([]Square, 0)
	moves = append(moves, b.getValidMovesHelper((*Square).up, Square{x, y})...)
	moves = append(moves, b.getValidMovesHelper((*Square).down, Square{x, y})...)
	moves = append(moves, b.getValidMovesHelper((*Square).left, Square{x, y})...)
	moves = append(moves, b.getValidMovesHelper((*Square).right, Square{x, y})...)
	return moves
}

func computeIncrement(src int, dst int) int {
	if src > dst {
		return -1
	}
	return 1
}

func (m *Move) decomposeMoveToSteps() []Step {
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
		b.Set(step1.Dst.X, step1.Dst.Y, stack(u1, udst))
		udst = b.pickup(step2.Dst.X, step2.Dst.Y)
		b.Set(step2.Dst.X, step2.Dst.Y, stack(u2, udst))
		return true, true
	}
	// Case 2: Steps are on adjacent squares, moving into one another
	// In this case, the stacking is applied, but neither piece will move
	if step1.Dst == step2.Src && step1.Src == step2.Dst {
		collision := stack(u1, u2)
		if collision.Exists && collision.Team == u1.Team {
			b.Set(step1.Src.X, step1.Src.Y, collision)
		} else if collision.Exists && collision.Team == u2.Team {
			b.Set(step2.Src.X, step2.Src.Y, collision)
		}
		return true, true
	}
	// Case 3: Two exclusive moves
	collision1 := b.Get(step1.Dst.X, step1.Dst.Y).Exists
	collision2 := b.Get(step2.Dst.X, step2.Dst.Y).Exists
	b.Set(step1.Dst.X, step1.Dst.Y, stack(u1, b.Get(step1.Dst.X, step1.Dst.Y)))
	b.Set(step2.Dst.X, step2.Dst.Y, stack(u2, b.Get(step2.Dst.X, step2.Dst.Y)))
	return collision1, collision2
}

func (b *Board) checkWinCondition() (bool, Team) {
	team1win := false
	team2win := false
	// Let rank 0 be team 1's "endzone"
	for i := 0; i < b.Cols; i++ {
		if b.Get(i, 0).Exists && b.Get(i, 0).Team == 1 {
			team1win = true
		}
	}
	// Let rank nRows be team 2's "endzone"
	for i := 0; i < b.Cols; i++ {
		if b.Get(i, b.Rows-1).Exists && b.Get(i, b.Rows-1).Team == 2 {
			team2win = true
		}
	}
	if team1win && team2win {
		// TODO: compare stack sizes
		return true, 0
	} else if team1win {
		return true, 1
	} else if team2win {
		return true, 2
	}
	return false, 0
}

func (b *Board) ResolveMove(move1 Move, move2 Move) (bool, Team) {
	// TODO: validate moves
	// logrus.WithFields(logrus.Fields{"Board": b}).Info("init")
	// defer logrus.WithFields(logrus.Fields{"Board": b}).Info("fini")
	steps1 := move1.decomposeMoveToSteps()
	steps2 := move2.decomposeMoveToSteps()

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
		winner, team := b.checkWinCondition()
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
