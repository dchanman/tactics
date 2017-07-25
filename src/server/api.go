package server

import (
	"github.com/dchanman/tactics/src/game"
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

// TacticsApi exposes game APIs to the client
type TacticsApi struct {
	id   uint64
	game *game.Game
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
	*result = TacticsApiResult{}
	return nil
}
