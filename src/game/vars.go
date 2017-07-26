package game

import "github.com/sirupsen/logrus"

var (
	log                  = logrus.WithField("pkg", "game")
	gameSubscriberBuffer = 10
	gameChatBuffer       = 10
)
