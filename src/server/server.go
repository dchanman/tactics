package server

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	maxMessageSize = 2048
)

// Server manages the communications with clients in order to manage games
type Server struct {
	games map[uint32]*game.Game
	maxid uint64
}

// NewServer instantiates a new server
func NewServer() *Server {
	games := make(map[uint32]*game.Game)
	return &Server{games: games, maxid: 1}
}

// CreateNewGame creates a game with a given ID.
// Returns error if game already exists
func (s *Server) CreateNewGame(gameid uint32, gameType game.GameType) error {
	if s.DoesGameIDExist(gameid) {
		return errors.New("game already exists")
	}
	s.games[gameid] = game.NewGame(gameType)
	return nil
}

// DoesGameIDExist returns true if a game ID is linked to an actual game
func (s *Server) DoesGameIDExist(gameid uint32) bool {
	_, ok := s.games[gameid]
	return ok
}

// GetGameIDs returns a list of all existing game IDs
func (s *Server) GetGameIds(req *http.Request, args *struct{}, result *struct {
	GameIDs []uint32 `json:"gameids"`
}) error {
	ids := make([]uint32, 0)
	for id := range s.games {
		ids = append(ids, id)
	}
	*result = struct {
		GameIDs []uint32 `json:"gameids"`
	}{ids}
	return nil
}

// CreateGame creates a new game and returns the game's ID
func (s *Server) CreateGame(req *http.Request, args *struct {
	GameType game.GameType `json:"gameType"`
}, result *struct {
	GameID uint32 `json:"gameid"`
}) error {
	log.WithFields(logrus.Fields{"gametype": args.GameType}).Info("Creating Game of type")
	var randID uint32
	for randID = generateRandomID(); s.DoesGameIDExist(randID); randID = generateRandomID() {
	}
	s.CreateNewGame(randID, args.GameType)
	*result = struct {
		GameID uint32 `json:"gameid"`
	}{randID}
	return nil
}

func generateRandomID() uint32 {
	return uint32(rand.Intn(999998) + 1)
}

func (s *Server) nextID() uint64 {
	ret := s.maxid
	s.maxid++
	return ret
}

// RegisterNewClient registers a new websocket connection with the server
func (s *Server) RegisterNewClient(gameid uint32, conn *websocket.Conn) error {
	if game, ok := s.games[gameid]; ok {
		api := NewTacticsApi(s.nextID(), conn)
		go api.subscribeToGame(game)
		go api.serveRPC()
		return nil
	}
	return errors.New("game ID not found")
}
