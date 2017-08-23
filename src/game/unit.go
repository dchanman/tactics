package game

import "fmt"

// Team is a number
type Team int8

// Unit is a basic unit in the game
type Unit struct {
	Team   Team `json:"team,omitempty"`
	Stack  int8 `json:"stack,omitempty"`
	Exists bool `json:"exists"`
}

func stack(u1 Unit, u2 Unit) Unit {
	if !u1.Exists {
		return u2
	}
	if !u2.Exists {
		return u1
	}
	if u1.Team == u2.Team {
		u1.Stack += u2.Stack
		return u1
	}
	if u1.Stack > u2.Stack {
		u1.Stack -= u2.Stack
		return u1
	}
	if u2.Stack > u1.Stack {
		u2.Stack -= u1.Stack
		return u2
	}
	return Unit{Exists: false}
}

func (u *Unit) getValidMoves(b *Board, sq Square) []Square {
	moves := make([]Square, 0)
	moves = append(moves, b.getLineInDirection((*Square).up, sq)...)
	moves = append(moves, b.getLineInDirection((*Square).down, sq)...)
	moves = append(moves, b.getLineInDirection((*Square).left, sq)...)
	moves = append(moves, b.getLineInDirection((*Square).right, sq)...)
	return moves
}

// String representation
func (u *Unit) String() string {
	if !u.Exists {
		return "."
	}
	return fmt.Sprintf("T%d:%d", u.Team, u.Stack)
}
