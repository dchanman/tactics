package server

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 2048
)

// Server manages the communications with clients in order to manage games
type Server struct {
	games map[uint32]*game.Game
	chats map[uint32]*game.Chat
	maxid uint64
}

// NewServer instantiates a new server
func NewServer() *Server {
	games := make(map[uint32]*game.Game)
	chats := make(map[uint32]*game.Chat)
	return &Server{games: games, chats: chats, maxid: 1}
}

func (s *Server) createNewGame(gameid uint32, gameType game.BoardType) error {
	if s.DoesGameIDExist(gameid) {
		return errors.New("game already exists")
	}
	s.games[gameid] = game.NewGame(gameType)
	s.chats[gameid] = game.NewChat()
	return nil
}

// DoesGameIDExist returns true if a game ID is linked to an actual game
func (s *Server) DoesGameIDExist(gameid uint32) bool {
	_, ok := s.games[gameid]
	return ok
}

// GetGameIds returns a list of all existing game IDs
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
	BoardType game.BoardType `json:"gameType"`
}, result *struct {
	GameID uint32 `json:"gameid"`
}) error {
	var randID uint32
	for randID = generateRandomID(); s.DoesGameIDExist(randID); randID = generateRandomID() {
	}
	s.createNewGame(randID, args.BoardType)
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
func (s *Server) RegisterNewClient(conn *websocket.Conn) error {
	api := NewTacticsApi(s.nextID(), conn, s)
	go api.serveRPC()
	return nil
}
