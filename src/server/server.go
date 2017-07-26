package server

import (
	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 2048
)

// Server manages the communications with clients in order to manage games
type Server struct {
	Game  *game.Game
	maxid uint64
}

// NewServer instantiates a new server
func NewServer() *Server {
	g := game.NewGame()
	g.B.Set(3, 4, game.Unit{Name: "hi", Exists: true})
	g.B.Set(5, 4, game.Unit{Name: "hi", Exists: true})
	return &Server{Game: g}
}

func (s *Server) nextID() uint64 {
	ret := s.maxid
	s.maxid++
	return ret
}

// RegisterNewClient registers a new websocket connection with the server
func (s *Server) RegisterNewClient(conn *websocket.Conn) {
	api := NewTacticsApi(s.nextID(), conn)
	go api.SubscribeToGame(s.Game)
	go api.ServeRPC()
}
