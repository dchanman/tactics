package game

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type chatMsg struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type Chat struct {
	in    chan chatMsg
	out   map[uint64](chan Notification)
	mutex sync.Mutex
}

func NewChat() *Chat {
	chat := Chat{
		in:  make(chan chatMsg, gameChatBuffer),
		out: make(map[uint64](chan Notification))}
	go chat.chatPump()
	return &chat
}

func (c *Chat) Send(id string, msg string) {
	c.in <- chatMsg{
		Sender:  id,
		Message: msg}
}

func (c *Chat) Subscribe(id uint64) chan Notification {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.out[id]; exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Duplicate ID")
		return nil
	}
	c.out[id] = make(chan Notification, gameChatBuffer)
	return c.out[id]
}

func (c *Chat) Unsubscribe(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.out[id]; !exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Unsubscribing invalid ID")
		return
	}
	delete(c.out, id)
}

func (c *Chat) chatPump() {
	for msg := range c.in {
		c.mutex.Lock()
		for _, ch := range c.out {
			notif := Notification{
				Method: "Game.Chat",
				Params: msg}
			ch <- notif
		}
		c.mutex.Unlock()
	}
}
