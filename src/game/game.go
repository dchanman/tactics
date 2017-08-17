package game

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	nMaxPlayers       = 2
	rowsLarge         = 10
	colsLarge         = 8
	rowsOfPiecesLarge = 3
	rowsSmall         = 8
	colsSmall         = 5
	rowsOfPiecesSmall = 2
)

// Game is the main game engine
type Game struct {
	gameType    GameType
	board       *Board
	subscribers map[uint64](chan *GameNotification)
	chat        chan GameChat
	movesQueue  chan gameMove
	history     []Turn
	completed   bool

	teamToPlayerID []uint64
	player1ready   bool
	player2ready   bool
}

// GameType determines the type of the game (board size, pieces, etc)
type GameType string

const (
	// Small 8x5 2 row configuration
	GameTypeSmall GameType = "small"
	// Large 10x8 3 row configuration
	GameTypeLarge = "large"
)

func (gt *GameType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	ret := GameType(s)
	if ret == GameTypeSmall || ret == GameTypeLarge {
		*gt = ret
		return nil
	}
	return errors.New("invalid game type")
}

type gameMove struct {
	move  Move
	team  Team
	reset bool
}

type Turn map[Team]Move

func NewTurn(move1 Move, move2 Move) Turn {
	moves := make(map[Team]Move)
	moves[1] = move1
	moves[2] = move2
	return Turn(moves)
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
	Board       *Board `json:"board"`
	History     []Turn `json:"history"`
	P1Available bool   `json:"p1available"`
	P2Available bool   `json:"p2available"`
	P1Ready     bool   `json:"p1ready"`
	P2Ready     bool   `json:"p2ready"`
}

func NewGame(gameType GameType) *Game {
	game := Game{
		gameType:       gameType,
		subscribers:    make(map[uint64](chan *GameNotification)),
		chat:           make(chan GameChat, gameChatBuffer),
		movesQueue:     make(chan gameMove, gameMoveBuffer),
		teamToPlayerID: make([]uint64, nMaxPlayers+1)}
	game.initGameState()
	go game.chatPump()
	go game.waitForMoves()
	return &game
}

func (g *Game) initGameState() {
	g.board = createGameBoard(g.gameType)
	g.history = make([]Turn, 0)
	g.completed = false
}

func (g *Game) ResetBoard() {
	g.initGameState()
	g.movesQueue <- gameMove{reset: true}
	g.PublishUpdate()
}

func createGameBoard(gameType GameType) *Board {
	var nCols int
	var nRows int
	var nRowsOfPieces int
	switch gameType {
	case GameTypeSmall:
		nCols = colsSmall
		nRows = rowsSmall
		nRowsOfPieces = rowsOfPiecesSmall
	case GameTypeLarge:
		nCols = colsLarge
		nRows = rowsLarge
		nRowsOfPieces = rowsOfPiecesLarge
	}
	b := NewBoard(nCols, nRows)
	// Add pieces
	for i := 0; i < nCols; i++ {
		for j := 0; j < nRowsOfPieces; j++ {
			b.Set(i, 1+j, Unit{Team: 2, Stack: 1, Exists: true})
			b.Set(i, nRows-2-j, Unit{Team: 1, Stack: 1, Exists: true})
		}
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
	var team Team
	var move1 Move
	var move2 Move
	g.player1ready = false
	g.player2ready = false
	for m := range g.movesQueue {
		winner := false
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
			winner, team = g.board.ResolveMove(move1, move2)
			g.history = append(g.history, NewTurn(move1, move2))
		}
		g.PublishUpdate()
		if winner {
			g.completed = true
			g.PublishVictory(team)
		}
	}
}

func (g *Game) getTeamForPlayerId(id uint64) Team {
	for team, pid := range g.teamToPlayerID {
		if pid == id {
			return Team(team)
		}
	}
	return 0
}

func (g *Game) CommitMove(id uint64, move Move) error {
	team := g.getTeamForPlayerId(id)
	if g.completed {
		return errors.New("Game is over")
	}
	if team == 0 {
		return errors.New("Player is not playing")
	}
	if !g.board.Get(move.Src.X, move.Src.Y).Exists || g.board.Get(move.Src.X, move.Src.Y).Team != team {
		return errors.New("Player is moving the wrong piece")
	}
	g.movesQueue <- gameMove{team: team, move: move}
	return nil
}

func (g *Game) GetValidMoves(id uint64, x int, y int) []Square {
	u := g.board.Get(x, y)
	if !g.completed && u.Exists && u.Team == g.getTeamForPlayerId(id) {
		return g.board.GetValidMoves(x, y)
	}
	return make([]Square, 0)
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
		Board:       g.board,
		History:     g.history,
		P1Available: g.teamToPlayerID[1] != 0,
		P2Available: g.teamToPlayerID[2] != 0,
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

func (g *Game) PublishVictory(team Team) {
	notif := GameNotification{
		Method: "Game.Over",
		Params: struct {
			Team Team `json:"team"`
		}{team}}
	for _, ch := range g.subscribers {
		ch <- &notif
	}
}

func (g *Game) GetPlayerReadyStatus() (bool, bool) {
	return g.player1ready, g.player2ready
}

func (g *Game) GetPlayerIds() (uint64, uint64) {
	return g.teamToPlayerID[1], g.teamToPlayerID[2]
}

func (g *Game) JoinGame(team int, id uint64) bool {
	if team > 0 && team <= nMaxPlayers {
		if g.teamToPlayerID[team] == 0 {
			g.teamToPlayerID[team] = id
			log.WithFields(logrus.Fields{"p1": g.teamToPlayerID[1], "p2": g.teamToPlayerID[2]}).Info("Joined")
			go g.PublishUpdate()
			return true
		}
	}
	log.WithFields(
		logrus.Fields{
			"pid":  id,
			"team": team,
			"p1":   g.teamToPlayerID[1],
			"p2":   g.teamToPlayerID[2]}).
		Error("Could not join")
	return false
}

func (g *Game) QuitGame(id uint64) {
	team := g.getTeamForPlayerId(id)
	if team > 0 {
		g.teamToPlayerID[team] = 0
		log.WithFields(logrus.Fields{"id": id}).Info("Quit")
		go g.PublishUpdate()
	}
}
