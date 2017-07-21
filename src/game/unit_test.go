package game

import "testing"

func TestUnitIsDead(t *testing.T) {
	u1 := Unit{Hp: 1}
	if u1.IsDead() {
		t.Error("Expected u1 to be alive!")
	}
	u2 := Unit{Hp: 0}
	if !u2.IsDead() {
		t.Error("Expected u2 to be dead!")
	}
}

func TestUnitAttack(t *testing.T) {
	u1 := Unit{
		Atk: 4,
	}
	u2 := Unit{
		Hp:  5,
		Def: 2,
	}
	u1.Attack(&u2)
	if u2.Hp != 3 {
		t.Error("Combat math mistake!")
	}
	u1.Attack(&u2)
	u1.Attack(&u2)
	if u2.Hp != 0 || !u2.IsDead() {
		t.Error("Combat math error: Expected 0, got", u2.Hp)
	}
}

func TestUnitAttackNoDamage(t *testing.T) {
	u1 := Unit{
		Atk: 4,
	}
	u2 := Unit{
		Hp:  5,
		Def: 4,
	}
	u1.Attack(&u2)
	if u2.Hp != 5 {
		t.Error("Combat math mistake!")
	}
	u2.Def = 5
	u1.Attack(&u2)
	if u2.Hp != 5 {
		t.Error("Combat math mistake!")
	}
}
