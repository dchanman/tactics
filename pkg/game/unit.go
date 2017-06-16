package game

// unit is a basic unit in the game
type unit struct {
	name   string
	class  string
	team   string
	hp     int8
	atk    int8
	def    int8
	mov    int8
	exists bool
}

func (u *unit) IsDead() bool {
	return (u.hp == 0)
}

func (u *unit) Attack(other *unit) {
	dmg := u.atk - other.def
	if dmg < 0 {
		dmg = 0
	}
	other.hp -= dmg
	if other.hp < 0 {
		other.hp = 0
	}
}
