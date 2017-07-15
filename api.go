package main

import "github.com/sirupsen/logrus"

type TacticsApiArgs struct {
	S string
}

type TacticsApiResult int

// TacticsApi exposes game APIs to the client
type TacticsApi struct {
	id   uint64
	game Game
}

func (api *TacticsApi) Hello(args *TacticsApiArgs, result *TacticsApiResult) error {
	log.WithFields(logrus.Fields{"args": args, "id": api.id}).Printf("Hello API called")
	*result = TacticsApiResult(0)
	return nil
}
