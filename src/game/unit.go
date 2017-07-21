package game

// unit is a basic unit in the game
type Unit struct {
	Name   string
	Class  string
	Team   string
	Hp     int8
	Atk    int8
	Def    int8
	Mov    int8
	Exists bool
}

func (u *Unit) IsDead() bool {
	return (u.Hp == 0)
}

func (u *Unit) Attack(other *Unit) {
	dmg := u.Atk - other.Def
	if dmg < 0 {
		dmg = 0
	}
	other.Hp -= dmg
	if other.Hp < 0 {
		other.Hp = 0
	}
}
