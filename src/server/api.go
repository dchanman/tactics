package server

import (
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/dchanman/tactics/src/game"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type TacticsApiArgs struct {
	X    int        `json:"x,omitempty"`
	Y    int        `json:"y,omitempty"`
	Unit *game.Unit `json:"unit,omitempty"`
}

type TacticsApiResult struct {
	Game *game.Game `json:"game,omitempty"`
}

type TacticsApiUpdate struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
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
	api := TacticsApi{id: id, gameFin: make(chan bool), client: &client}
	return &api
}

func (api *TacticsApi) SubscribeToGame(game *game.Game) {
	api.game = game
	ch := game.Subscribe(api.id)
	defer game.Unsubscribe(api.id)
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
	}()
	rpcserver := rpc.NewServer()
	rpcserver.Register(api)
	rpcserver.ServeCodec(jsonrpc.NewServerCodec(api.client))
}

func (api *TacticsApi) Hello(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Hello API called")
	*result = TacticsApiResult{}
	return nil
}

func (api *TacticsApi) GetGame(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Getting Game")
	log.WithFields(logrus.Fields{"game": api.game}).Printf("Game")
	*result = TacticsApiResult{Game: api.game}
	return nil
}

func (api *TacticsApi) AddUnit(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Adding unit")
	api.game.B.Set(args.X, args.Y, *args.Unit)
	go api.game.PublishUpdate()
	*result = TacticsApiResult{}
	return nil
}
