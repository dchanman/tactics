package game

import "testing"

func TestGetSet(t *testing.T) {
	b := newBoard(3, 4)
	u := unit{name: "test", exists: true}
	tu := b.get(0, 0)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(0, 0, u)
	tu = b.get(0, 0)
	if !tu.exists || tu.name != "test" {
		t.Error("Unexpected unit")
	}

	tu = b.get(2, 3)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(2, 3, u)
	tu = b.get(2, 3)
	if !tu.exists || tu.name != "test" {
		t.Error("Unexpected unit")
	}
}

func TestGetSetOutOfBounds(t *testing.T) {
	b := newBoard(3, 4)
	u := unit{name: "test", exists: true}
	tu := b.get(3, 4)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(3, 4, u)
	tu = b.get(3, 4)
	if tu.exists {
		t.Error("Unexpected unit")
	}

	tu = b.get(4, 3)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(4, 3, u)
	tu = b.get(4, 3)
	if tu.exists {
		t.Error("Unexpected unit")
	}

	tu = b.get(0, 4)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(0, 4, u)
	tu = b.get(0, 4)
	if tu.exists {
		t.Error("Unexpected unit")
	}

	tu = b.get(3, 0)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(3, 0, u)
	tu = b.get(3, 0)
	if tu.exists {
		t.Error("Unexpected unit")
	}

	tu = b.get(-1, -1)
	if tu.exists {
		t.Error("Unexpected unit")
	}
	b.set(-1, -1, u)
	tu = b.get(-1, -1)
	if tu.exists {
		t.Error("Unexpected unit")
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
	u := unit{mov: 0}
	b.set(2, 2, u)
	moves := b.getValidMoves(2, 2)
	if len(moves) != 0 {
		t.Error("Unexpected valid moves ", moves)
	}
}
func TestGetValidMoves1(t *testing.T) {
	b := newBoard(5, 5)
	u := unit{mov: 1}
	b.set(2, 2, u)
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
	b := newBoard(5, 5)
	u := unit{mov: 2}
	b.set(2, 2, u)
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
	b := newBoard(5, 5)
	u := unit{mov: 3}
	b.set(2, 2, u)
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
	b := newBoard(5, 5)
	u := unit{mov: 4}
	b.set(2, 2, u)
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
	b := newBoard(5, 5)
	u := unit{mov: 127}
	b.set(2, 2, u)
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
