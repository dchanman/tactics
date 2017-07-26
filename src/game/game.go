package game

import "github.com/sirupsen/logrus"

const ()

// Game is the main game engine
type Game struct {
	B           Board `json:"board,omitempty"`
	subscribers map[uint64](chan *GameNotification)
	chat        chan string
}

type GameNotification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewGame() *Game {
	game := Game{
		B:           NewBoard(10, 10),
		subscribers: make(map[uint64](chan *GameNotification)),
		chat:        make(chan string, gameChatBuffer)}
	go game.chatPump()
	return &game
}

func (g *Game) chatPump() {
	for msg := range g.chat {
		notif := GameNotification{
			Method: "Game.Chat",
			Params: struct {
				Message string `json:"message"`
			}{Message: msg}}
		for _, ch := range g.subscribers {
			ch <- &notif
		}
	}
}

func (g *Game) SendChat(msg string) {
	g.chat <- msg
}

func (g *Game) Subscribe(id uint64) chan *GameNotification {
	if _, exists := g.subscribers[id]; exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Duplicate ID")
		return nil
	}
	g.subscribers[id] = make(chan *GameNotification, gameSubscriberBuffer)
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
	notif := GameNotification{
		Method: "Game.Update",
		Params: struct {
			Game *Game `json:"game"`
		}{Game: g}}
	for _, ch := range g.subscribers {
		ch <- &notif
	}
}
