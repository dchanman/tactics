package game

// unit is a basic unit in the game
type Unit struct {
	Name   string `json:"name,omitempty"`
	Team   int8   `json:"team,omitempty"`
	Stack  int8   `json:"stack,omitempty"`
	Exists bool   `json:"exists"`
}

func stack(u1 Unit, u2 Unit) Unit {
	if !u1.Exists {
		return u2
	}
	if !u2.Exists {
		return u1
	}
	if u1.Stack > u2.Stack {
		if u1.Team != u2.Team {
			u1.Stack -= u2.Stack
		} else {
			u1.Stack += u2.Stack
		}
		return u1
	}
	if u1.Stack < u2.Stack {
		if u1.Team != u2.Team {
			u2.Stack -= u1.Stack
		} else {
			u2.Stack += u1.Stack
		}
		return u2
	}
	return Unit{Exists: false}
}
