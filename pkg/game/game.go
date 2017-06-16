package game

import "github.com/sirupsen/logrus"

var (
	log = logrus.WithField("pkg", "game")
)

// Game is the main game engine
type Game struct {
	name string
}
