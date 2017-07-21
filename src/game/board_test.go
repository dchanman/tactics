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
	u := Unit{Exists: true, Mov: 0}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	if len(moves) != 0 {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMoves1(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true, Mov: 1}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 4 ||
		!verify[square{3, 2}] ||
		!verify[square{1, 2}] ||
		!verify[square{2, 1}] ||
		!verify[square{2, 3}] {
		t.Error("Unexpected valid moves ", moves, verify)
	}
}
func TestGetValidMoves2(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true, Mov: 2}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 12 ||
		!verify[square{4, 2}] ||
		!verify[square{3, 1}] ||
		!verify[square{3, 2}] ||
		!verify[square{3, 3}] ||
		!verify[square{2, 0}] ||
		!verify[square{2, 1}] ||
		!verify[square{2, 3}] ||
		!verify[square{2, 4}] ||
		!verify[square{1, 1}] ||
		!verify[square{1, 2}] ||
		!verify[square{1, 3}] ||
		!verify[square{0, 2}] {
		t.Error("Unexpected valid moves ", moves, verify)
	}
}
func TestGetValidMoves3(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true, Mov: 3}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 20 ||
		!verify[square{4, 1}] ||
		!verify[square{4, 2}] ||
		!verify[square{4, 3}] ||
		!verify[square{3, 0}] ||
		!verify[square{3, 1}] ||
		!verify[square{3, 2}] ||
		!verify[square{3, 3}] ||
		!verify[square{3, 4}] ||
		!verify[square{2, 0}] ||
		!verify[square{2, 1}] ||
		!verify[square{2, 3}] ||
		!verify[square{2, 4}] ||
		!verify[square{1, 0}] ||
		!verify[square{1, 1}] ||
		!verify[square{1, 2}] ||
		!verify[square{1, 3}] ||
		!verify[square{1, 4}] ||
		!verify[square{0, 1}] ||
		!verify[square{0, 2}] ||
		!verify[square{0, 3}] {
		t.Error("Unexpected valid moves ", len(moves), moves)
	}
}
func TestGetValidMoves4(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true, Mov: 4}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 24 ||
		!verify[square{4, 0}] ||
		!verify[square{4, 1}] ||
		!verify[square{4, 2}] ||
		!verify[square{4, 3}] ||
		!verify[square{4, 4}] ||
		!verify[square{3, 0}] ||
		!verify[square{3, 1}] ||
		!verify[square{3, 2}] ||
		!verify[square{3, 3}] ||
		!verify[square{3, 4}] ||
		!verify[square{2, 0}] ||
		!verify[square{2, 1}] ||
		!verify[square{2, 3}] ||
		!verify[square{2, 4}] ||
		!verify[square{1, 0}] ||
		!verify[square{1, 1}] ||
		!verify[square{1, 2}] ||
		!verify[square{1, 3}] ||
		!verify[square{1, 4}] ||
		!verify[square{0, 0}] ||
		!verify[square{0, 1}] ||
		!verify[square{0, 2}] ||
		!verify[square{0, 3}] ||
		!verify[square{0, 4}] {
		t.Error("Unexpected valid moves ", moves, verify)
	}
}
func TestGetValidMoves127(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true, Mov: 127}
	b.Set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 24 ||
		!verify[square{4, 0}] ||
		!verify[square{4, 1}] ||
		!verify[square{4, 2}] ||
		!verify[square{4, 3}] ||
		!verify[square{4, 4}] ||
		!verify[square{3, 0}] ||
		!verify[square{3, 1}] ||
		!verify[square{3, 2}] ||
		!verify[square{3, 3}] ||
		!verify[square{3, 4}] ||
		!verify[square{2, 0}] ||
		!verify[square{2, 1}] ||
		!verify[square{2, 3}] ||
		!verify[square{2, 4}] ||
		!verify[square{1, 0}] ||
		!verify[square{1, 1}] ||
		!verify[square{1, 2}] ||
		!verify[square{1, 3}] ||
		!verify[square{1, 4}] ||
		!verify[square{0, 0}] ||
		!verify[square{0, 1}] ||
		!verify[square{0, 2}] ||
		!verify[square{0, 3}] ||
		!verify[square{0, 4}] {
		t.Error("Unexpected valid moves ", moves, verify)
	}
}
func TestGetValidMovesEnemyPieces(t *testing.T) {
	b := NewBoard(5, 5)
	u1 := Unit{Exists: true, Team: "a", Mov: 2}
	u2 := Unit{Exists: true, Team: "b"}
	u3 := Unit{Exists: true, Team: "a"}
	b.Set(2, 2, u1)
	b.Set(2, 1, u2)
	b.Set(2, 3, u3)
	moves := b.getValidMoves(2, 2)
	verify := make(map[square]bool, len(moves))
	for i := range moves {
		verify[moves[i]] = true
	}
	if len(moves) != 9 ||
		!verify[square{4, 2}] ||
		!verify[square{3, 1}] ||
		!verify[square{3, 2}] ||
		!verify[square{3, 3}] ||
		!verify[square{2, 4}] ||
		!verify[square{1, 1}] ||
		!verify[square{1, 2}] ||
		!verify[square{1, 3}] ||
		!verify[square{0, 2}] {
		t.Error("Unexpected valid moves ", len(moves), moves)
	}
}
func TestToJSON(t *testing.T) {
	b := NewBoard(5, 5)
	u1 := Unit{Exists: true, Team: "a", Mov: 2}
	u2 := Unit{Exists: true, Team: "b"}
	u3 := Unit{Exists: true, Team: "a"}
	b.Set(2, 2, u1)
	b.Set(2, 1, u2)
	b.Set(2, 3, u3)
	b.ToJSON()
}
