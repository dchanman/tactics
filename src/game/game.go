package game

import "github.com/sirupsen/logrus"

const (
	nRows = 10
	nCols = 8
	// nRows = 8
	// nCols = 5
)

// Game is the main game engine
type Game struct {
	B           *Board `json:"board,omitempty"`
	subscribers map[uint64](chan *GameNotification)
	chat        chan GameChat
	movesQueue  chan gameMove

	player1id    uint64
	player2id    uint64
	player1ready bool
	player2ready bool
}

type gameMove struct {
	move  Move
	team  Team
	reset bool
}

type GameChat struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type GameNotification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type GameInformation struct {
	Game        *Game `json:"game"`
	P1Available bool  `json:"p1available"`
	P2Available bool  `json:"p2available"`
	P1Ready     bool  `json:"p1ready"`
	P2Ready     bool  `json:"p2ready"`
}

func NewGame() *Game {
	game := Game{
		B:           createGameBoard(),
		subscribers: make(map[uint64](chan *GameNotification)),
		chat:        make(chan GameChat, gameChatBuffer),
		movesQueue:  make(chan gameMove, gameMoveBuffer),
		player1id:   0,
		player2id:   0}
	go game.chatPump()
	go game.waitForMoves()
	return &game
}

func createGameBoard() *Board {
	b := NewBoard(nCols, nRows)
	// Add pieces
	for i := 0; i < nCols; i++ {
		b.Set(i, 1, Unit{Team: 2, Stack: 1, Exists: true})
		b.Set(i, 2, Unit{Team: 2, Stack: 1, Exists: true})
		b.Set(i, 3, Unit{Team: 2, Stack: 1, Exists: true})
		b.Set(i, nRows-2, Unit{Team: 1, Stack: 1, Exists: true})
		b.Set(i, nRows-3, Unit{Team: 1, Stack: 1, Exists: true})
		b.Set(i, nRows-4, Unit{Team: 1, Stack: 1, Exists: true})
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
	g.player1ready = false
	g.player2ready = false
	for m := range g.movesQueue {
		if m.reset {
			g.player1ready = false
			g.player2ready = false
		} else if m.team == 1 {
			log.WithFields(logrus.Fields{"move": m.move}).Info("Player 1 ready")
			g.player1ready = true
			move1 = m.move
		} else if m.team == 2 {
			log.WithFields(logrus.Fields{"move": m.move}).Info("Player 2 ready")
			g.player2ready = true
			move2 = m.move
		}
		if g.player1ready && g.player2ready {
			g.player1ready = false
			g.player2ready = false
			g.B.ResolveMove(move1, move2)
		}
		g.PublishUpdate()
	}
}

func (g *Game) CommitMove(team Team, move Move) {
	g.movesQueue <- gameMove{team: team, move: move}
}

func (g *Game) ResetBoard() {
	g.B = createGameBoard()
	g.movesQueue <- gameMove{reset: true}
	g.PublishUpdate()
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

func (g *Game) GetGameInformation() GameInformation {
	return GameInformation{
		Game:        g,
		P1Available: !(g.player1id == 0),
		P2Available: !(g.player2id == 0),
		P1Ready:     g.player1ready,
		P2Ready:     g.player2ready}
}

func (g *Game) PublishUpdate() {
	notif := GameNotification{
		Method: "Game.Update",
		Params: g.GetGameInformation()}
	for _, ch := range g.subscribers {
		ch <- &notif
	}
}

func (g *Game) GetPlayerReadyStatus() (bool, bool) {
	return g.player1ready, g.player2ready
}

func (g *Game) GetPlayerIds() (uint64, uint64) {
	return g.player1id, g.player2id
}

func (g *Game) JoinGame(playerNumber int, id uint64) bool {
	if playerNumber == 1 {
		if g.player1id == 0 {
			g.player1id = id
			log.WithFields(logrus.Fields{"p1": g.player1id, "p2": g.player2id}).Info("Joined")
			go g.PublishUpdate()
			return true
		}
	} else if playerNumber == 2 {
		if g.player2id == 0 {
			g.player2id = id
			log.WithFields(logrus.Fields{"p1": g.player1id, "p2": g.player2id}).Info("Joined")
			go g.PublishUpdate()
			return true
		}
	}
	log.WithFields(
		logrus.Fields{
			"pid":          id,
			"playerNumber": playerNumber,
			"p1":           g.player1id,
			"p2":           g.player2id}).
		Error("Could not join")
	return false
}

func (g *Game) QuitGame(id uint64) {
	if g.player1id == id {
		g.player1id = 0
		log.WithFields(logrus.Fields{"id": id}).Info("Quit")
		go g.PublishUpdate()
	}
	if g.player2id == id {
		g.player2id = 0
		log.WithFields(logrus.Fields{"id": id}).Info("Quit")
		go g.PublishUpdate()
	}
}
