package game

// unit is a basic unit in the game
type Unit struct {
	Name   string `json:"name,omitempty"`
	Class  string `json:"class,omitempty"`
	Team   string `json:"team,omitempty"`
	Hp     int8   `json:"hp,omitempty"`
	Atk    int8   `json:"atk,omitempty"`
	Def    int8   `json:"def,omitempty"`
	Mov    int8   `json:"mov,omitempty"`
	Exists bool   `json:"exists"`
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
