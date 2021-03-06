package server

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type chatMsg struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// Chat provides chatroom functionality
type Chat struct {
	in    chan chatMsg
	out   map[uint64](chan Notification)
	mutex sync.Mutex
}

// NewChat instantiates a new Chat
func NewChat() *Chat {
	chat := Chat{
		in:  make(chan chatMsg, gameChatBuffer),
		out: make(map[uint64](chan Notification))}
	go chat.chatPump()
	return &chat
}

// Send allows a client to send a message
func (c *Chat) Send(id string, msg string) {
	c.in <- chatMsg{
		Sender:  id,
		Message: msg}
}

// Subscribe allows a client to receive chat notifications
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

// Unsubscribe allows a client to stop receiving chat notifications
func (c *Chat) Unsubscribe(id uint64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.out[id]; !exists {
		log.WithFields(logrus.Fields{"id": id}).Error("Unsubscribing invalid ID")
		return
	}
	delete(c.out, id)
}

func (c *Chat) broadcastNotification(notif Notification) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, ch := range c.out {
		ch <- notif
	}
}

func (c *Chat) announceJoin(name string) {
	notif := Notification{
		Method: "Chat.Join",
		Params: name}
	c.broadcastNotification(notif)
}

func (c *Chat) announceNameChange(from string, to string) {
	notif := Notification{
		Method: "Chat.NameChange",
		Params: struct {
			From string `json:"from"`
			To   string `json:"to"`
		}{from, to}}
	c.broadcastNotification(notif)
}

func (c *Chat) chatPump() {
	for msg := range c.in {
		notif := Notification{
			Method: "Chat.Message",
			Params: msg}
		c.broadcastNotification(notif)
	}
}
