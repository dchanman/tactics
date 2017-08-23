package server

import (
	"errors"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"

	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	TacticsApiErrorNoID   = errors.New("ID not found")
	TacticsApiErrorNoGame = errors.New("no registered game")
	TacticsApiErrorNoChat = errors.New("no registered chat")
)

type TacticsApiRole string

const (
	TacticsApiRoleObserver TacticsApiRole = "Spectator"
	TacticsApiRolePlayer                  = "Player"
)

// TacticsApi exposes game APIs to the client
type TacticsApi struct {
	server  *Server
	id      uint64
	game    *game.Game
	chat    *game.Chat
	gameFin map[game.Subscribable](chan bool)
	done    bool
	client  *Client
}

func NewTacticsApi(id uint64, conn *websocket.Conn, s *Server) *TacticsApi {
	client := Client{conn: conn}
	api := TacticsApi{
		id:      id,
		gameFin: make(map[game.Subscribable](chan bool)),
		client:  &client,
		server:  s}
	return &api
}

func (api *TacticsApi) subscribeAndServe(s game.Subscribable) {
	ch := s.Subscribe(api.id)
	defer s.Unsubscribe(api.id)
	fin := make(chan bool)
	api.gameFin[s] = fin
	for {
		select {
		case update := <-ch:
			api.client.WriteJSON(update)
		case <-fin:
			log.WithFields(logrus.Fields{"id": api.id}).Info("Terminating pump")
			return
		}
	}
}

func (api *TacticsApi) serveRPC() {
	defer func() {
		log.Info("Done Serving")
		if r := recover(); r != nil {
			log.WithFields(logrus.Fields{"r": r}).Info("Recovered")
		}
		err := api.client.conn.Close()
		if err != nil {
			log.Error("Close: ", err)
		}
		for _, fin := range api.gameFin {
			fin <- true
		}
		if api.game != nil {
			api.game.QuitGame(api.id)
		}
	}()
	rpcserver := rpc.NewServer()
	rpcserver.Register(api)
	rpcserver.ServeCodec(jsonrpc.NewServerCodec(api.client))
}

func (api *TacticsApi) Heartbeat(args *struct{}, result *struct{}) error {
	return nil
}

func (api *TacticsApi) GetGame(args *struct{}, result *game.Information) error {
	if api.game == nil {
		return TacticsApiErrorNoGame
	}
	*result = api.game.GetInformation()
	return nil
}

func (api *TacticsApi) SendChat(args *struct {
	Message string `json:"message"`
}, result *struct{}) error {
	if api.chat == nil {
		return TacticsApiErrorNoChat
	}
	api.chat.Send(strconv.FormatUint(api.id, 10), args.Message)
	return nil
}

func (api *TacticsApi) CommitMove(args *struct {
	FromX int `json:"fromX"`
	FromY int `json:"fromY"`
	ToX   int `json:"toX"`
	ToY   int `json:"toY"`
}, result *struct{}) error {
	if api.game == nil {
		return TacticsApiErrorNoGame
	}
	return api.game.CommitMove(api.id, game.Square{X: args.FromX, Y: args.FromY}, game.Square{X: args.ToX, Y: args.ToY})
}

func (api *TacticsApi) ResetBoard(args *struct{}, result *struct{}) error {
	if api.game == nil {
		return TacticsApiErrorNoGame
	}
	api.game.ResetBoard()
	return nil
}

func (api *TacticsApi) JoinGame(args *struct {
	PlayerNumber int `json:"playerNumber"`
}, result *struct{}) error {
	if api.game == nil {
		return TacticsApiErrorNoGame
	}
	joinedResult := api.game.JoinGame(args.PlayerNumber, api.id)
	if !joinedResult {
		return errors.New("Could not join game")
	}
	return nil
}

func (api *TacticsApi) GetRole(args *struct{}, result *struct {
	Role TacticsApiRole `json:"role"`
	Team game.Team      `json:"team"`
}) error {
	if api.game == nil {
		return TacticsApiErrorNoGame
	}
	p1id, p2id := api.game.GetPlayerIDs()
	role := TacticsApiRoleObserver
	team := game.Team(0)
	if api.id == p1id {
		role = TacticsApiRolePlayer
		team = 1
	}
	if api.id == p2id {
		role = TacticsApiRolePlayer
		team = 2
	}

	*result = struct {
		Role TacticsApiRole `json:"role"`
		Team game.Team      `json:"team"`
	}{
		Role: role,
		Team: team}
	return nil
}

func (api *TacticsApi) SubscribeGame(args *struct {
	ID uint32 `json:"id"`
}, result *struct{}) error {
	// TODO: threadsafety
	log.WithFields(logrus.Fields{"id": args.ID}).Info("SubG")
	game, ok := api.server.games[args.ID]
	if !ok {
		return TacticsApiErrorNoID
	}
	if api.game != nil {
		fin := api.gameFin[api.game]
		fin <- true
		delete(api.gameFin, api.game)
	}
	api.game = game
	go api.subscribeAndServe(game)
	return nil
}

func (api *TacticsApi) SubscribeChat(args *struct {
	ID uint32 `json:"id"`
}, result *struct{}) error {
	log.WithFields(logrus.Fields{"id": args.ID}).Info("SubC")
	chat, ok := api.server.chats[args.ID]
	if !ok {
		return TacticsApiErrorNoID
	}
	if api.chat != nil {
		fin := api.gameFin[api.chat]
		fin <- true
		delete(api.gameFin, api.chat)
	}
	api.chat = chat
	go api.subscribeAndServe(chat)
	return nil
}
