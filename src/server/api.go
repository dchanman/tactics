package server

import (
	"github.com/dchanman/tactics/src/game"
	"github.com/sirupsen/logrus"
)

type TacticsApiArgs struct {
	S string
}

type TacticsApiResult struct {
	Game *game.Board `json:"game,omitempty"`
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
	*result = TacticsApiResult{Game: &api.game.B}
	return nil
}
