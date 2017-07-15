package main

import (
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	maxMessageSize = 2048
)

// Server manages the communications with clients in order to manage games
type Server struct {
	Game  Game
	maxid uint64
}

// NewServer instantiates a new server
func NewServer() *Server {
	g := NewGame()
	g.b.set(3, 4, unit{Name: "hi", Exists: true})
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
	api := TacticsApi{id: s.nextID(), game: s.Game}
	rpcserver := rpc.NewServer()
	rpcserver.Register(&api)
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
		}()
		rpcserver.ServeCodec(jsonrpc.NewServerCodec(&c))
	}()
}
