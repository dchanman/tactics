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

type TacticsApiArgs struct {
	X       int        `json:"x,omitempty"`
	Y       int        `json:"y,omitempty"`
	Unit    *game.Unit `json:"unit,omitempty"`
	Message string     `json:"message,omitempty"`
}

type TacticsApiResult struct {
	Game       *game.Game    `json:"game,omitempty"`
	ValidMoves []game.Square `json:"validMoves,omitempty"`
}

type TacticsApiUpdate struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type TacticsApiRole string

const (
	TacticsApiRoleObserver TacticsApiRole = "Spectator"
	TacticsApiRolePlayer                  = "Player"
)

type TacticsApiPlayerStatus string

const (
	TacticsApiPlayerStatusUnavailable TacticsApiPlayerStatus = "Unavailable"
	TacticsApiPlayerStatusThinking                           = "Thinking of move..."
	TacticsApiPlayerStatusCommitted                          = "Move committed"
)

type TacticsApiStatus struct {
	Role     TacticsApiRole         `json:"role"`
	P1Status TacticsApiPlayerStatus `json:"p1status"`
	P2Status TacticsApiPlayerStatus `json:"p2status"`
}

// TacticsApi exposes game APIs to the client
type TacticsApi struct {
	id      uint64
	game    *game.Game
	gameFin chan bool
	client  *Client
}

func NewTacticsApi(id uint64, conn *websocket.Conn) *TacticsApi {
	client := Client{conn: conn}
	api := TacticsApi{
		id:      id,
		gameFin: make(chan bool),
		client:  &client}
	return &api
}

func (api *TacticsApi) SubscribeToGame(g *game.Game) {
	api.game = g
	ch := g.Subscribe(api.id)
	defer g.Unsubscribe(api.id)
	for {
		select {
		case update := <-ch:
			log.WithFields(logrus.Fields{"id": api.id}).Info("Updated!")
			api.client.WriteJSON(update)
		case <-api.gameFin:
			log.WithFields(logrus.Fields{"id": api.id}).Info("Terminating pump")
			return
		}
	}
}

func (api *TacticsApi) ServeRPC() {
	defer func() {
		log.Info("Done Serving")
		if r := recover(); r != nil {
			log.WithFields(logrus.Fields{"r": r}).Info("Recovered")
		}
		err := api.client.conn.Close()
		if err != nil {
			log.Error("Close: ", err)
		}
		api.gameFin <- true
		api.game.QuitGame(api.id)
	}()
	rpcserver := rpc.NewServer()
	rpcserver.Register(api)
	rpcserver.ServeCodec(jsonrpc.NewServerCodec(api.client))
}

func (api *TacticsApi) Heartbeat(args *struct{}, result *struct{}) error {
	log.WithFields(logrus.Fields{"id": api.id}).Printf("Heartbeat")
	return nil
}

func (api *TacticsApi) GetGame(args *TacticsApiArgs, result *game.GameInformation) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Getting Game")
	log.WithFields(logrus.Fields{"game": api.game}).Printf("Game")
	*result = api.game.GetGameInformation()
	return nil
}

func (api *TacticsApi) AddUnit(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Adding unit")
	api.game.B.Set(args.X, args.Y, *args.Unit)
	go api.game.PublishUpdate()
	*result = TacticsApiResult{}
	return nil
}

func (api *TacticsApi) SendChat(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Sending Chat")
	api.game.SendChat(strconv.FormatUint(api.id, 10), args.Message)
	return nil
}

func (api *TacticsApi) GetValidMoves(args *TacticsApiArgs, result *TacticsApiResult) error {
	// TODO: Eventually this will be done clientside
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Getting moves")
	*result = TacticsApiResult{ValidMoves: api.game.B.GetValidMoves(args.X, args.Y)}
	return nil
}

func (api *TacticsApi) CommitMove(args *struct {
	FromX int `json:"fromX"`
	FromY int `json:"fromY"`
	ToX   int `json:"toX"`
	ToY   int `json:"toY"`
}, result *struct{}) error {
	// TODO: validate move
	// TODO: assign teams
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Committing move")
	team := game.Team(api.id%2) + 1
	move := game.Move{
		Src: game.Square{X: args.FromX, Y: args.FromY},
		Dst: game.Square{X: args.ToX, Y: args.ToY}}
	api.game.CommitMove(team, move)
	return nil
}

func (api *TacticsApi) ResetBoard(args *struct{}, result *struct{}) error {
	log.WithFields(logrus.Fields{"id": api.id}).Printf("Resetting board")
	api.game.ResetBoard()
	return nil
}

func (api *TacticsApi) JoinGame(args *struct {
	PlayerNumber int `json:"playerNumber"`
}, result *struct{}) error {
	joinedResult := api.game.JoinGame(args.PlayerNumber, api.id)
	if !joinedResult {
		return errors.New("Could not join game")
	}
	return nil
}

func (api *TacticsApi) GetRole(args *struct{}, result *struct {
	Role TacticsApiRole `json:"role"`
}) error {
	p1id, p2id := api.game.GetPlayerIds()
	var role TacticsApiRole
	role = TacticsApiRoleObserver
	if api.id == p1id || api.id == p2id {
		role = TacticsApiRolePlayer
	}

	*result = struct {
		Role TacticsApiRole `json:"role"`
	}{
		Role: role}
	return nil
}
