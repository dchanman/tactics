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
