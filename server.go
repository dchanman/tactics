package main

import (
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	maxMessageSize = 2048
)

type TacticsApiArgs struct {
	S string
}

type TacticsApiResult int

// TacticsApi exposes game APIs to the client
type TacticsApi struct {
	id uint64
}

func (api *TacticsApi) Hello(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Hello API called")
	*result = TacticsApiResult(0)
	return nil
}

// Server manages the communications with clients in order to manage games
type server struct {
	clients []client
	Game    Game
	maxid   uint64
}

type client struct {
	conn   *websocket.Conn
	reader io.Reader
	writer io.WriteCloser
}

func (c *client) Read(p []byte) (n int, err error) {
	if c.reader == nil {
		_, c.reader, err = c.conn.NextReader()
		if err != nil {
			n = 0
			return
		}
	}
	for n = 0; n < len(p); {
		var bytes int
		bytes, err = c.reader.Read(p[n:])
		n += bytes
		if err == io.EOF {
			c.reader = nil
			break
		}
		if err != nil {
			break
		}
	}
	return
}

func (c *client) Write(p []byte) (n int, err error) {
	if c.writer == nil {
		c.writer, err = c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			n = 0
			return
		}
	}
	for n = 0; n < len(p); {
		var bytes int
		bytes, err = c.writer.Write(p)
		n += bytes
		if err != nil {
			break
		}
	}
	if err != nil || n == len(p) {
		err = c.Close()
	}
	return
}

func (c *client) Close() (err error) {
	if c.writer != nil {
		err = c.writer.Close()
		c.writer = nil
	}
	return
}

func NewServer() *server {
	g := NewGame()
	g.b.set(3, 4, unit{Name: "hi", Exists: true})
	return &server{Game: g}
}

func (s *server) nextID() uint64 {
	ret := s.maxid
	s.maxid++
	return ret
}

func (s *server) registerNewClient(conn *websocket.Conn) {
	c := client{conn: conn}
	api := TacticsApi{id: s.nextID()}
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
