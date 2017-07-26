package server

import (
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	maxMessageSize = 2048
)

// Server manages the communications with clients in order to manage games
type Server struct {
	Game  game.Game
	maxid uint64
}

type PushMsg struct {
	Method string           `json:"method"`
	Params TacticsApiResult `json:"params"`
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
	c := Client{conn: conn}
	api := TacticsApi{id: s.nextID(), game: &s.Game}
	rpcserver := rpc.NewServer()
	rpcserver.Register(&api)
	// TODO: Refactor
	fin := make(chan bool)
	go func() {
		ch := s.Game.Subscribe(api.id)
		for {
			select {
			case <-ch:
				log.WithFields(logrus.Fields{"id": api.id}).Info("Updated!")
				c.WriteJSON(PushMsg{Method: "TacticsApi.Update", Params: TacticsApiResult{Game: &s.Game}})
			case <-fin:
				log.WithFields(logrus.Fields{"id": api.id}).Info("Terminating pump")
				s.Game.Unsubscribe(api.id)
				return
			}
		}
	}()
	go func() {
		defer func() {
			log.Info("Done Serving")
			if r := recover(); r != nil {
				log.WithFields(logrus.Fields{"r": r}).Info("Recovered")
			}
			err := conn.Close()
			if err != nil {
				log.Error("Close: ", err)
			}
			fin <- true
		}()
		rpcserver.ServeCodec(jsonrpc.NewServerCodec(&c))
	}()
}
