package server

import "github.com/sirupsen/logrus"

var (
	log                  = logrus.WithField("pkg", "server")
	gameSubscriberBuffer = 10
	gameChatBuffer       = 10
	gameMoveBuffer       = 2
)
