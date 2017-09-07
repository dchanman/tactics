package server

import (
	"errors"
	"net/rpc"
	"net/rpc/jsonrpc"
	"regexp"
	"strconv"

	"github.com/dchanman/tactics/src/game"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	errNoID        = errors.New("ID not found")
	errNoGame      = errors.New("no registered game")
	errNoChat      = errors.New("no registered chat")
	errBadUsername = errors.New("invalid username")
	usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9_.-]{1,20}$")
)

// Role is a user's role in a game
type Role string

const (
	roleObserver Role = "Spectator"
	rolePlayer        = "Player"
)

// TacticsAPI exposes game APIs to the client
type TacticsAPI struct {
	server  *Server
	id      uint64
	name    string
	game    *Game
	chat    *Chat
	gameFin map[Subscribable](chan bool)
	done    bool
	client  *Client
}

func newTacticsAPI(id uint64, conn *websocket.Conn, s *Server) *TacticsAPI {
	client := Client{conn: conn}
	api := TacticsAPI{
		id:      id,
		name:    "Guest-" + strconv.FormatUint(id, 10),
		gameFin: make(map[Subscribable](chan bool)),
		client:  &client,
		server:  s}
	return &api
}

func (api *TacticsAPI) subscribeAndServe(s Subscribable) {
	ch := s.Subscribe(api.id)
	defer s.Unsubscribe(api.id)
	fin := make(chan bool)
	api.gameFin[s] = fin
	for {
		select {
		case update := <-ch:
			api.client.WriteJSON(update)
		case <-fin:
			return
		}
	}
}

func (api *TacticsAPI) serveRPC() {
	defer func() {
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

func (api *TacticsAPI) Heartbeat(args *struct{}, result *struct{}) error {
	return nil
}

func (api *TacticsAPI) GetGame(args *struct{}, result *Information) error {
	if api.game == nil {
		return errNoGame
	}
	*result = api.game.GetInformation()
	return nil
}

func (api *TacticsAPI) SendChat(args *struct {
	Message string `json:"message"`
}, result *struct{}) error {
	if api.chat == nil {
		return errNoChat
	}
	api.chat.Send(api.name, args.Message)
	return nil
}

func (api *TacticsAPI) CommitMove(args *struct {
	FromX int `json:"fromX"`
	FromY int `json:"fromY"`
	ToX   int `json:"toX"`
	ToY   int `json:"toY"`
}, result *struct{}) error {
	if api.game == nil {
		return errNoGame
	}
	from := game.Square{X: args.FromX, Y: args.FromY}
	to := game.Square{X: args.ToX, Y: args.ToY}
	return api.game.CommitMove(api.id, from, to)
}

func (api *TacticsAPI) ResetBoard(args *struct{}, result *struct{}) error {
	if api.game == nil {
		return errNoGame
	}
	api.game.ResetBoard()
	return nil
}

func (api *TacticsAPI) JoinGame(args *struct {
	PlayerNumber int `json:"playerNumber"`
}, result *struct{}) error {
	if api.game == nil {
		return errNoGame
	}
	joinedResult := api.game.JoinGame(args.PlayerNumber, api.id)
	if !joinedResult {
		return errors.New("Could not join game")
	}
	return nil
}

func (api *TacticsAPI) GetRole(args *struct{}, result *struct {
	Role Role      `json:"role"`
	Team game.Team `json:"team"`
}) error {
	if api.game == nil {
		return errNoGame
	}
	p1id, p2id := api.game.GetPlayerIDs()
	role := roleObserver
	team := game.Team(0)
	if api.id == p1id {
		role = rolePlayer
		team = 1
	}
	if api.id == p2id {
		role = rolePlayer
		team = 2
	}

	*result = struct {
		Role Role      `json:"role"`
		Team game.Team `json:"team"`
	}{
		Role: role,
		Team: team}
	return nil
}

func (api *TacticsAPI) SubscribeGame(args *struct {
	ID uint32 `json:"id"`
}, result *struct{}) error {
	// TODO: threadsafety
	game, ok := api.server.games[args.ID]
	if !ok {
		return errNoID
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

func (api *TacticsAPI) SubscribeChat(args *struct {
	ID uint32 `json:"id"`
}, result *struct{}) error {
	chat, ok := api.server.chats[args.ID]
	if !ok {
		return errNoID
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

func (api *TacticsAPI) SetChatName(args *struct {
	Name string `json:"name"`
}, result *struct{}) error {
	log.WithFields(logrus.Fields{"name": args.Name, "match": usernameRegexp.MatchString(args.Name)}).Info("Set Chat")
	if !usernameRegexp.MatchString(args.Name) {
		return errBadUsername
	}
	api.name = args.Name
	return nil
}

func (api *TacticsAPI) GetChatName(args *struct{}, result *struct {
	Name string `json:"name"`
}) error {
	*result = struct {
		Name string `json:"name"`
	}{api.name}
	return nil
}
