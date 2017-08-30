package server

import (
	"encoding/json"
	"errors"
	"game"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	nMaxPlayers       = 2
	rowsLarge         = 10
	colsLarge         = 7
	rowsOfPiecesLarge = 3
	rowsSmall         = 8
	colsSmall         = 5
	rowsOfPiecesSmall = 2
)

var (
	errInvalidBoardType = errors.New("invalid game type")
	errGameOver         = errors.New("Game is over")
	errPlayerNotPlaying = errors.New("Player is not playing")
	errInvalidSquare    = errors.New("invalid square")
	errPlayerWrongTeam  = errors.New("Player is moving the wrong piece")
)

// Game is the main game engine
type Game struct {
	gameType     BoardType
	board        *game.Board
	initialBoard *game.Board
	subscribers  map[uint64](chan Notification)
	movesQueue   chan gameMove
	history      []Turn
	completed    bool

	teamToPlayerID      []uint64
	teamToPlayerIDMutex sync.Mutex
	player1ready        bool
	player2ready        bool
}

// BoardType determines the type of the game (board size, pieces, etc)
type BoardType string

const (
	// BoardTypeSmall is a 8x5 2 row configuration
	BoardTypeSmall BoardType = "small"
	// BoardTypeLarge is a 10x8 3 row configuration
	BoardTypeLarge = "large"
)

// UnmarshalJSON validates BoardTypes
func (gt *BoardType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	ret := BoardType(s)
	if ret == BoardTypeSmall || ret == BoardTypeLarge {
		*gt = ret
		return nil
	}
	return errInvalidBoardType
}

type gameMove struct {
	move  game.Move
	team  game.Team
	reset bool
}

// Turn is a set of moves made by all players
type Turn map[game.Team]game.Move

// NewTurn constructs a new Turn
func NewTurn(move1 game.Move, move2 game.Move) Turn {
	moves := make(Turn)
	moves[1] = move1
	moves[2] = move2
	return moves
}

// Information contains player information
type Status struct {
	P1Available bool `json:"p1available"`
	P2Available bool `json:"p2available"`
	P1Ready     bool `json:"p1ready"`
	P2Ready     bool `json:"p2ready"`
}
type Information struct {
	Status
	InitialBoard *game.Board `json:"board"`
	History      []Turn      `json:"history"`
}

// NewGame constructs a new Game
func NewGame(gameType BoardType) *Game {
	game := Game{
		gameType:       gameType,
		initialBoard:   createGameBoard(gameType),
		subscribers:    make(map[uint64](chan Notification)),
		movesQueue:     make(chan gameMove, gameMoveBuffer),
		teamToPlayerID: make([]uint64, nMaxPlayers+1)}
	game.initGameState()
	go game.waitForMoves()
	return &game
}

func (g *Game) initGameState() {
	// TODO: just copy g.initialGameBoard
	g.board = createGameBoard(g.gameType)
	g.history = make([]Turn, 0)
	g.completed = false
}

// ResetBoard resets the board
func (g *Game) ResetBoard() {
	g.initGameState()
	g.movesQueue <- gameMove{reset: true}
	g.publishUpdate()
}

func createGameBoard(gameType BoardType) *game.Board {
	var nCols int
	var nRows int
	var nRowsOfPieces int
	switch gameType {
	case BoardTypeSmall:
		nCols = colsSmall
		nRows = rowsSmall
		nRowsOfPieces = rowsOfPiecesSmall
	case BoardTypeLarge:
		nCols = colsLarge
		nRows = rowsLarge
		nRowsOfPieces = rowsOfPiecesLarge
	}
	b := game.NewBoard(nCols, nRows)
	// Add pieces
	for i := 0; i < nCols; i++ {
		for j := 0; j < nRowsOfPieces; j++ {
			b.Set(i, 1+j, game.Unit{Team: 2, Stack: 1, Exists: true})
			b.Set(i, nRows-2-j, game.Unit{Team: 1, Stack: 1, Exists: true})
		}
	}
	return &b
}

func (g *Game) waitForMoves() {
	var move1 game.Move
	var move2 game.Move
	var resolution game.Resolution
	g.player1ready = false
	g.player2ready = false
	for m := range g.movesQueue {
		winner := false
		if m.reset {
			g.player1ready = false
			g.player2ready = false
		} else if m.team == 1 {
			g.player1ready = true
			move1 = m.move
		} else if m.team == 2 {
			g.player2ready = true
			move2 = m.move
		}
		if g.player1ready && g.player2ready {
			g.player1ready = false
			g.player2ready = false
			resolution = g.board.ResolveMove(move1, move2)
			winner = resolution.Winner
			g.history = append(g.history, NewTurn(move1, move2))
			g.publishTurn()
		}
		g.publishStatus()
		if winner {
			g.completed = true
			g.publishVictory(resolution.Team)
		}
	}
}

func (g *Game) getTeamForPlayerID(id uint64) game.Team {
	for team, pid := range g.teamToPlayerID {
		if pid == id {
			return game.Team(team)
		}
	}
	return 0
}

// CommitMove is called from a client to commit a move on behalf of a player
func (g *Game) CommitMove(id uint64, src game.Square, dst game.Square) error {
	team := g.getTeamForPlayerID(id)
	if g.completed {
		return errGameOver
	}
	if team == 0 {
		return errPlayerNotPlaying
	}
	if !g.board.IsValid(src.X, src.Y) || !g.board.IsValid(dst.X, dst.Y) {
		return errInvalidSquare
	}
	if !g.board.Get(src.X, src.Y).Exists || g.board.Get(src.X, src.Y).Team != team {
		return errPlayerWrongTeam
	}
	move := game.Move{Src: src, Dst: dst}
	g.movesQueue <- gameMove{team: team, move: move}
	return nil
}

func (g *Game) GetValidMoves(id uint64, x int, y int) []game.Square {
	u := g.board.Get(x, y)
	if !g.completed && u.Exists && u.Team == g.getTeamForPlayerID(id) {
		return g.board.GetValidMoves(x, y)
	}
	return make([]game.Square, 0)
}

// Subscribe implements Subscriber interface
func (g *Game) Subscribe(id uint64) chan Notification {
	if _, exists := g.subscribers[id]; exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Duplicate ID")
		return nil
	}
	g.subscribers[id] = make(chan Notification, gameSubscriberBuffer)
	return g.subscribers[id]
}

// Unsubscribe implements Subscriber interface
func (g *Game) Unsubscribe(id uint64) {
	if _, exists := g.subscribers[id]; !exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Unsubscribing invalid ID")
		return
	}
	delete(g.subscribers, id)
}

func (g *Game) getStatus() Status {
	return Status{
		P1Available: g.teamToPlayerID[1] != 0,
		P2Available: g.teamToPlayerID[2] != 0,
		P1Ready:     g.player1ready,
		P2Ready:     g.player2ready}
}

// GetInformation returns the current status of the game
func (g *Game) GetInformation() Information {
	return Information{
		InitialBoard: g.initialBoard,
		History:      g.history,
		Status:       g.getStatus()}
}

func (g *Game) publishUpdate() {
	notif := Notification{
		Method: "Game.Update",
		Params: g.GetInformation()}
	for _, ch := range g.subscribers {
		ch <- notif
	}
}

func (g *Game) publishStatus() {
	notif := Notification{
		Method: "Game.Status",
		Params: g.getStatus()}
	for _, ch := range g.subscribers {
		ch <- notif
	}
}

func (g *Game) publishTurn() {
	notif := Notification{
		Method: "Game.Turn",
		Params: g.history}
	for _, ch := range g.subscribers {
		ch <- notif
	}
}

func (g *Game) publishVictory(team game.Team) {
	notif := Notification{
		Method: "Game.Over",
		Params: struct {
			Team game.Team `json:"team"`
		}{team}}
	for _, ch := range g.subscribers {
		ch <- notif
	}
}

func (g *Game) getPlayerReadyStatus() (bool, bool) {
	return g.player1ready, g.player2ready
}

// GetPlayerIDs returns the IDs of the players currently playing the game
func (g *Game) GetPlayerIDs() (uint64, uint64) {
	return g.teamToPlayerID[1], g.teamToPlayerID[2]
}

// JoinGame sets an ID as a player currently playing the game
func (g *Game) JoinGame(team int, id uint64) bool {
	g.teamToPlayerIDMutex.Lock()
	defer g.teamToPlayerIDMutex.Unlock()
	if team > 0 && team <= nMaxPlayers {
		if g.teamToPlayerID[team] == 0 {
			g.teamToPlayerID[team] = id
			g.publishStatus()
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

// QuitGame allows a player currently playing the game to quit
func (g *Game) QuitGame(id uint64) {
	g.teamToPlayerIDMutex.Lock()
	defer g.teamToPlayerIDMutex.Unlock()
	team := g.getTeamForPlayerID(id)
	if team > 0 {
		g.teamToPlayerID[team] = 0
		go g.publishUpdate()
	}
}
