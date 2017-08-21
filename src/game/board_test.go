package game

import "testing"

func TestGetSet(t *testing.T) {
	b := NewBoard(3, 4)
	u := Unit{Name: "test", Exists: true}
	tu := b.Get(0, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(0, 0, u)
	tu = b.Get(0, 0)
	if !tu.Exists || tu.Name != "test" {
		t.Error("Unexpected Unit")
	}

	tu = b.Get(2, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(2, 3, u)
	tu = b.Get(2, 3)
	if !tu.Exists || tu.Name != "test" {
		t.Error("Unexpected Unit")
	}
}

func TestGetSetOutOfBounds(t *testing.T) {
	b := NewBoard(3, 4)
	u := Unit{Name: "test", Exists: true}
	tu := b.Get(3, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(3, 4, u)
	tu = b.Get(3, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.Get(4, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(4, 3, u)
	tu = b.Get(4, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.Get(0, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(0, 4, u)
	tu = b.Get(0, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.Get(3, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(3, 0, u)
	tu = b.Get(3, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.Get(-1, -1)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.Set(-1, -1, u)
	tu = b.Get(-1, -1)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
}

func TestGetValidMovesNil(t *testing.T) {
	b := NewBoard(5, 5)
	moves := b.getValidMoves(2, 2)
	if len(moves) != 0 {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMoves0(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[Square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 8 ||
		!verify[Square{2, 1}] ||
		!verify[Square{2, 0}] ||
		!verify[Square{2, 3}] ||
		!verify[Square{2, 4}] ||
		!verify[Square{1, 2}] ||
		!verify[Square{0, 2}] ||
		!verify[Square{3, 2}] ||
		!verify[Square{4, 2}] {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMovesEnemyPieces(t *testing.T) {
	b := NewBoard(5, 5)
	u1 := Unit{Exists: true, Team: 1}
	u2 := Unit{Exists: true, Team: 2}
	u3 := Unit{Exists: true, Team: 1}
	b.Set(2, 2, u1)
	b.Set(2, 1, u2)
	b.Set(2, 3, u3)
	moves := b.getValidMoves(2, 2)
	verify := make(map[Square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 8 ||
		!verify[Square{2, 1}] ||
		!verify[Square{2, 0}] ||
		!verify[Square{2, 3}] ||
		!verify[Square{2, 4}] ||
		!verify[Square{1, 2}] ||
		!verify[Square{0, 2}] ||
		!verify[Square{3, 2}] ||
		!verify[Square{4, 2}] {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestDecomposeMoveToSteps(t *testing.T) {
	var result []Step
	var move Move

	move = Move{Src: Square{0, 0}, Dst: Square{0, 3}}
	result = move.decomposeMoveToSteps()
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 0}, Dst: Square{0, 1}}) ||
		(result[1] != Step{Src: Square{0, 1}, Dst: Square{0, 2}}) ||
		(result[2] != Step{Src: Square{0, 2}, Dst: Square{0, 3}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{0, 3}, Dst: Square{0, 0}}
	result = move.decomposeMoveToSteps()
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 3}, Dst: Square{0, 2}}) ||
		(result[1] != Step{Src: Square{0, 2}, Dst: Square{0, 1}}) ||
		(result[2] != Step{Src: Square{0, 1}, Dst: Square{0, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{0, 0}, Dst: Square{3, 0}}
	result = move.decomposeMoveToSteps()
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 0}, Dst: Square{1, 0}}) ||
		(result[1] != Step{Src: Square{1, 0}, Dst: Square{2, 0}}) ||
		(result[2] != Step{Src: Square{2, 0}, Dst: Square{3, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{3, 0}, Dst: Square{0, 0}}
	result = move.decomposeMoveToSteps()
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{3, 0}, Dst: Square{2, 0}}) ||
		(result[1] != Step{Src: Square{2, 0}, Dst: Square{1, 0}}) ||
		(result[2] != Step{Src: Square{1, 0}, Dst: Square{0, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{1, 1}, Dst: Square{1, 1}}
	result = move.decomposeMoveToSteps()
	if len(result) != 0 {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{1, 1}, Dst: Square{2, 2}}
	result = move.decomposeMoveToSteps()
	if len(result) != 0 {
		t.Error("Unexpected steps ", result)
	}
}
func TestResolveStep(t *testing.T) {
	var b Board
	var s1 Step
	var s2 Step

	// Case 1: Moves converging to a single square
	// 	Subcase: Two equally stacked pieces converging
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(1, 0))
	}
	if b.Get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	// 	Subcase: Two differently stacked pieces converging
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 0, Unit{Stack: 2, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(1, 0))
	}
	if (b.Get(0, 0) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 3, Team: 1, Exists: true})
	b.Set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(1, 0))
	}
	if (b.Get(0, 0) != Unit{Stack: 2, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	// Subcase: Larger stack colliding, while a defensive stack occurs
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 0, Unit{Stack: 2, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(1, 0))
	}
	if b.Get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}

	// Case 2: Adjacent pieces moving into one another
	// 	Subcase: Two equally stacked pieces converging
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(1, 0))
	}
	if b.Get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	// 	Subcase: Two differently stacked pieces converging
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 2, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.Get(1, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(1, 0))
	}
	if (b.Get(0, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 2, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 4, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 1))
	}
	if (b.Get(0, 0) != Unit{Stack: 2, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}

	// Case 3: Independent moves
	// 	Subcase: pieces "following" one another
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 2}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.Get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	if (b.Get(0, 1) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 1))
	}
	if (b.Get(0, 2) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 2))
	}
	// 	Subcase: pieces "dodging"
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.Get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.Get(0, 0))
	}
	if (b.Get(0, 1) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(0, 1))
	}
	if (b.Get(1, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(1, 1))
	}
	// 	Subcase: completely independent
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{1, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(0, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(0, 0))
	}
	if (b.Get(1, 0) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(1, 0))
	}
	if (b.Get(1, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(1, 1))
	}
	// Subcase: completely independent friendly stacking
	b = NewBoard(3, 3)
	b.Set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{1, 0}}
	b.resolveStep(s1, s2)
	if b.Get(0, 1).Exists || b.Get(0, 0).Exists {
		t.Error("Unexpected pieces: ", b.Get(0, 1), " and ", b.Get(0, 0))
	}
	if (b.Get(1, 0) != Unit{Stack: 2, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(1, 0))
	}
	if (b.Get(1, 1) != Unit{Stack: 2, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.Get(1, 1))
	}

}
func TestResolveMoveWinConditions(t *testing.T) {
	var b Board
	var m1 Move
	var m2 Move
	var winner bool
	var team Team
	cols := 3
	rows := 5

	// Test race victory: team 1 arrives first
	b = NewBoard(cols, rows)
	b.Set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 2}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 1 {
		t.Error("Moves resolved incorrectly")
	}

	// Test race victory: draw
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 0 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side move short
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 2}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 1 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 1}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side friendly collision
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 2, Unit{Stack: 2, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 2, Team: 1, Exists: true})
	b.Set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side enemy collision
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 2, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side collision, which removes a blockade allowing other side win
	b = NewBoard(cols, rows)
	b.Set(1, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(1, 3, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(1, 4, Unit{Stack: 1, Team: 1, Exists: true})
	m1 = Move{Src: Square{1, 4}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{1, 1}, Dst: Square{1, 4}}
	winner, team = b.ResolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test both sides collision, nobody wins
	b = NewBoard(cols, rows)
	b.Set(1, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(0, 3, Unit{Stack: 1, Team: 1, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, _ = b.ResolveMove(m1, m2)
	if winner {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test collision, make sure pieces stop
	b = NewBoard(cols, rows)
	b.Set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.Set(1, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.Set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 2}}
	winner, _ = b.ResolveMove(m1, m2)
	if winner {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	if b.Get(1, 3).Exists ||
		!(b.Get(1, 2).Exists && b.Get(1, 2).Stack == 2) ||
		!(b.Get(1, 1).Exists && b.Get(1, 1).Stack == 1) {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
}
