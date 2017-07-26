package game

import "github.com/sirupsen/logrus"

// Game is the main game engine
type Game struct {
	B           Board `json:"board,omitempty"`
	subscribers map[uint64](chan string)
}

func NewGame() Game {
	return Game{B: NewBoard(10, 10), subscribers: make(map[uint64](chan string))}
}

func (g *Game) Subscribe(id uint64) chan string {
	if _, exists := g.subscribers[id]; exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Duplicate ID")
		return nil
	}
	g.subscribers[id] = make(chan string, gameSubscriberBuffer)
	return g.subscribers[id]
}

func (g *Game) Unsubscribe(id uint64) {
	if _, exists := g.subscribers[id]; !exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Unsubscribing invalid ID")
		return
	}
	delete(g.subscribers, id)
}

func (g *Game) PublishUpdate() {
	for _, ch := range g.subscribers {
		ch <- "update"
	}
}
