package game

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

type Step Move

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

func (b *Board) GetValidMoves(x int, y int) []Square {
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

func (b *Board) resolveStep(step1 Step, step2 Step) bool {
	u1 := b.pickup(step1.Src.X, step1.Src.Y)
	u2 := b.pickup(step2.Src.X, step2.Src.Y)
	// Case 1: Both steps converge onto the same square
	if step1.Dst == step2.Dst {
		b.Set(step1.Dst.X, step1.Dst.Y, stack(u1, u2))
		return true
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
		return true
	}
	// Case 3: Two exclusive moves
	b.Set(step1.Dst.X, step1.Dst.Y, stack(u1, b.Get(step1.Dst.X, step1.Dst.Y)))
	b.Set(step2.Dst.X, step2.Dst.Y, stack(u2, b.Get(step2.Dst.X, step2.Dst.Y)))
	return false
}
