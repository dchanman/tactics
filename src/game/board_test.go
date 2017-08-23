package game

import "testing"

func TestGetSet(t *testing.T) {
	b := newBoard(3, 4)
	u := Unit{Exists: true}
	tu := b.get(0, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(0, 0, u)
	tu = b.get(0, 0)
	if !tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.get(2, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(2, 3, u)
	tu = b.get(2, 3)
	if !tu.Exists {
		t.Error("Unexpected Unit")
	}
}

func TestGetSetOutOfBounds(t *testing.T) {
	b := newBoard(3, 4)
	u := Unit{Exists: true}
	tu := b.get(3, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(3, 4, u)
	tu = b.get(3, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.get(4, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(4, 3, u)
	tu = b.get(4, 3)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.get(0, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(0, 4, u)
	tu = b.get(0, 4)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.get(3, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(3, 0, u)
	tu = b.get(3, 0)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}

	tu = b.get(-1, -1)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
	b.set(-1, -1, u)
	tu = b.get(-1, -1)
	if tu.Exists {
		t.Error("Unexpected Unit")
	}
}

func TestSetNUnits(t *testing.T) {
	b := newBoard(2, 2)
	if b.nUnits != 0 {
		t.Error("Unexpected nUnits", b.nUnits)
	}
	// Place 1 piece
	b.set(0, 0, Unit{Exists: true, Stack: 1})
	if b.nUnits != 1 {
		t.Error("Unexpected nUnits", b.nUnits)
	}
	// Place another piece
	b.set(0, 1, Unit{Exists: true, Stack: 2})
	if b.nUnits != 3 {
		t.Error("Unexpected nUnits", b.nUnits)
	}
	// Replace another piece
	b.set(0, 0, Unit{Exists: true, Stack: 3})
	if b.nUnits != 5 {
		t.Error("Unexpected nUnits", b.nUnits)
	}
	// Remove all pieces
	b.set(0, 0, Unit{})
	b.set(0, 1, Unit{})
	if b.nUnits != 0 {
		t.Error("Unexpected nUnits", b.nUnits)
	}
}

func TestGetValidMovesNil(t *testing.T) {
	b := newBoard(5, 5)
	moves := b.getValidMoves(2, 2)
	if len(moves) != 0 {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMoves0(t *testing.T) {
	b := newBoard(5, 5)
	u := Unit{Exists: true}
	b.set(2, 2, u)
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
	b := newBoard(5, 5)
	u1 := Unit{Exists: true, Team: 1}
	u2 := Unit{Exists: true, Team: 2}
	u3 := Unit{Exists: true, Team: 1}
	b.set(2, 2, u1)
	b.set(2, 1, u2)
	b.set(2, 3, u3)
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
	result = decomposeMoveToSteps(move)
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 0}, Dst: Square{0, 1}}) ||
		(result[1] != Step{Src: Square{0, 1}, Dst: Square{0, 2}}) ||
		(result[2] != Step{Src: Square{0, 2}, Dst: Square{0, 3}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{0, 3}, Dst: Square{0, 0}}
	result = decomposeMoveToSteps(move)
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 3}, Dst: Square{0, 2}}) ||
		(result[1] != Step{Src: Square{0, 2}, Dst: Square{0, 1}}) ||
		(result[2] != Step{Src: Square{0, 1}, Dst: Square{0, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{0, 0}, Dst: Square{3, 0}}
	result = decomposeMoveToSteps(move)
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{0, 0}, Dst: Square{1, 0}}) ||
		(result[1] != Step{Src: Square{1, 0}, Dst: Square{2, 0}}) ||
		(result[2] != Step{Src: Square{2, 0}, Dst: Square{3, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{3, 0}, Dst: Square{0, 0}}
	result = decomposeMoveToSteps(move)
	if len(result) != 3 ||
		(result[0] != Step{Src: Square{3, 0}, Dst: Square{2, 0}}) ||
		(result[1] != Step{Src: Square{2, 0}, Dst: Square{1, 0}}) ||
		(result[2] != Step{Src: Square{1, 0}, Dst: Square{0, 0}}) {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{1, 1}, Dst: Square{1, 1}}
	result = decomposeMoveToSteps(move)
	if len(result) != 0 {
		t.Error("Unexpected steps ", result)
	}

	move = Move{Src: Square{1, 1}, Dst: Square{2, 2}}
	result = decomposeMoveToSteps(move)
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
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(1, 0))
	}
	if b.get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	// 	Subcase: Two differently stacked pieces converging
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 0, Unit{Stack: 2, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(1, 0))
	}
	if (b.get(0, 0) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 3, Team: 1, Exists: true})
	b.set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(1, 0))
	}
	if (b.get(0, 0) != Unit{Stack: 2, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	// Subcase: Larger stack colliding, while a defensive stack occurs
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 0, Unit{Stack: 2, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{1, 0}, Dst: Square{0, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(1, 0))
	}
	if b.get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}

	// Case 2: Adjacent pieces moving into one another
	// 	Subcase: Two equally stacked pieces converging
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(1, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(1, 0))
	}
	if b.get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	// 	Subcase: Two differently stacked pieces converging
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 2, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.get(1, 0).Exists {
		t.Error("Unexpected piece: ", b.get(1, 0))
	}
	if (b.get(0, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 2, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 4, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 0}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists {
		t.Error("Unexpected piece: ", b.get(0, 1))
	}
	if (b.get(0, 0) != Unit{Stack: 2, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}

	// Case 3: Independent moves
	// 	Subcase: pieces "following" one another
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{0, 2}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	if (b.get(0, 1) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 1))
	}
	if (b.get(0, 2) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 2))
	}
	// 	Subcase: pieces "dodging"
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{0, 1}}
	b.resolveStep(s1, s2)
	if b.get(0, 0).Exists {
		t.Error("Unexpected piece: ", b.get(0, 0))
	}
	if (b.get(0, 1) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(0, 1))
	}
	if (b.get(1, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(1, 1))
	}
	// 	Subcase: completely independent
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{1, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(0, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(0, 0))
	}
	if (b.get(1, 0) != Unit{Stack: 1, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(1, 0))
	}
	if (b.get(1, 1) != Unit{Stack: 1, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(1, 1))
	}
	// Subcase: completely independent friendly stacking
	b = newBoard(3, 3)
	b.set(0, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 1, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 0, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(1, 0, Unit{Stack: 1, Team: 2, Exists: true})
	s1 = Step{Src: Square{0, 1}, Dst: Square{1, 1}}
	s2 = Step{Src: Square{0, 0}, Dst: Square{1, 0}}
	b.resolveStep(s1, s2)
	if b.get(0, 1).Exists || b.get(0, 0).Exists {
		t.Error("Unexpected pieces: ", b.get(0, 1), " and ", b.get(0, 0))
	}
	if (b.get(1, 0) != Unit{Stack: 2, Team: 2, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(1, 0))
	}
	if (b.get(1, 1) != Unit{Stack: 2, Team: 1, Exists: true}) {
		t.Error("Unexpected piece: ", b.get(1, 1))
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
	b = newBoard(cols, rows)
	b.set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 2}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 1 {
		t.Error("Moves resolved incorrectly")
	}

	// Test race victory: draw
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 0 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test race victory: more pieces on team 1
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 2, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 1 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test race victory: more pieces on team 2
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 3, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 5}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side move short
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 2}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 1 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 1}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side friendly collision
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 2, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 2, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side victory, other side enemy collision
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test one side collision, which removes a blockade allowing other side win
	b = newBoard(cols, rows)
	b.set(1, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(1, 3, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(1, 4, Unit{Stack: 1, Team: 1, Exists: true})
	m1 = Move{Src: Square{1, 4}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{1, 1}, Dst: Square{1, 4}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 2 {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test collision, make sure pieces stop
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 1, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 2}}
	winner, _ = b.resolveMove(m1, m2)
	if winner {
		t.Error("Moves resolved incorrectly. Winner: ", winner)
	}
	if b.get(1, 3).Exists ||
		!(b.get(1, 2).Exists && b.get(1, 2).Stack == 2) ||
		!(b.get(1, 1).Exists && b.get(1, 1).Stack == 1) {
		t.Error("Moves resolved incorrectly. Team: ", team)
	}

	// Test collision, no pieces remaining, draw
	b = newBoard(cols, rows)
	b.set(1, 3, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(1, 2, Unit{Stack: 1, Team: 2, Exists: true})
	b.set(0, 2, Unit{Stack: 1, Team: 1, Exists: true})
	b.set(0, 1, Unit{Stack: 1, Team: 2, Exists: true})
	m1 = Move{Src: Square{1, 3}, Dst: Square{1, 0}}
	m2 = Move{Src: Square{0, 1}, Dst: Square{0, 3}}
	winner, team = b.resolveMove(m1, m2)
	if !winner || team != 0 {
		t.Error("Moves resolved incorrectly. Winner: ", winner, " Team: ", team, b)
	}
	if b.nUnits > 0 {
		t.Error("Moves resolved incorrectly. ", b.nUnits, " units leftover.", b)
	}
}
