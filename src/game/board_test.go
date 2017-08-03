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
	moves := b.GetValidMoves(2, 2)
	if len(moves) != 0 {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMoves0(t *testing.T) {
	b := NewBoard(5, 5)
	u := Unit{Exists: true}
	b.Set(2, 2, u)
	moves := b.GetValidMoves(2, 2)
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
	moves := b.GetValidMoves(2, 2)
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
