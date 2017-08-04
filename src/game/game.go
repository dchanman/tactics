package game

import "github.com/sirupsen/logrus"

const (
	nRows = 8
	nCols = 5
)

// Game is the main game engine
type Game struct {
	B           *Board `json:"board,omitempty"`
	subscribers map[uint64](chan *GameNotification)
	chat        chan GameChat
	movesQueue  chan gameMove
}

type gameMove struct {
	move Move
	team Team
}

type GameChat struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type GameNotification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewGame() *Game {
	game := Game{
		B:           createGameBoard(),
		subscribers: make(map[uint64](chan *GameNotification)),
		chat:        make(chan GameChat, gameChatBuffer),
		movesQueue:  make(chan gameMove, gameMoveBuffer)}
	go game.chatPump()
	go game.waitForMoves()
	return &game
}

func createGameBoard() *Board {
	b := NewBoard(nCols, nRows)
	// Add pieces
	for i := 0; i < nCols; i++ {
		b.Set(i, 1, Unit{Name: "X", Team: 1, Stack: 1, Exists: true})
		b.Set(i, 2, Unit{Name: "X", Team: 1, Stack: 1, Exists: true})
		b.Set(i, nRows-2, Unit{Name: "O", Team: 2, Stack: 1, Exists: true})
		b.Set(i, nRows-3, Unit{Name: "O", Team: 2, Stack: 1, Exists: true})
	}
	return &b
}

func (g *Game) chatPump() {
	for msg := range g.chat {
		notif := GameNotification{
			Method: "Game.Chat",
			Params: msg}
		for _, ch := range g.subscribers {
			ch <- &notif
		}
	}
}

func (g *Game) waitForMoves() {
	var move1 Move
	var move2 Move
	ready1 := false
	ready2 := false
	for m := range g.movesQueue {
		if m.team == 1 {
			log.WithFields(logrus.Fields{"move": m.move}).Info("Player 1 ready")
			ready1 = true
			move1 = m.move
		} else if m.team == 2 {
			log.WithFields(logrus.Fields{"move": m.move}).Info("Player 2 ready")
			ready2 = true
			move2 = m.move
		}
		if ready1 && ready2 {
			ready1 = false
			ready2 = false
			g.B.ResolveMove(move1, move2)
			g.PublishUpdate()
		}
	}
}

func (g *Game) CommitMove(team Team, move Move) {
	g.movesQueue <- gameMove{team: team, move: move}
}

func (g *Game) SendChat(sender string, msg string) {
	g.chat <- GameChat{
		Sender:  sender,
		Message: msg}
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
