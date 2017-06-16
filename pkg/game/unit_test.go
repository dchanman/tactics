package game

import "testing"

func TestUnitIsDead(t *testing.T) {
	u1 := unit{hp: 1}
	if u1.IsDead() {
		t.Error("Expected u1 to be alive!")
	}
	u2 := unit{hp: 0}
	if !u2.IsDead() {
		t.Error("Expected u2 to be dead!")
	}
}

func TestUnitAttack(t *testing.T) {
	u1 := unit{
		atk: 4,
	}
	u2 := unit{
		hp:  5,
		def: 2,
	}
	u1.Attack(&u2)
	if u2.hp != 3 {
		t.Error("Combat math mistake!")
	}
	u1.Attack(&u2)
	u1.Attack(&u2)
	if u2.hp != 0 || !u2.IsDead() {
		t.Error("Combat math error: Expected 0, got", u2.hp)
	}
}

func TestUnitAttackNoDamage(t *testing.T) {
	u1 := unit{
		atk: 4,
	}
	u2 := unit{
		hp:  5,
		def: 4,
	}
	u1.Attack(&u2)
	if u2.hp != 5 {
		t.Error("Combat math mistake!")
	}
	u2.def = 5
	u1.Attack(&u2)
	if u2.hp != 5 {
		t.Error("Combat math mistake!")
	}
}
