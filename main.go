package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	maxMessageSize = 2048
)

var (
	log      = logrus.WithField("pkg", "tactics")
	upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

type Client struct {
	conn *websocket.Conn
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, bytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error(err)
			}
			break
		}
		// TODO: safely convert bytes to string?
		msg := string(bytes)
		log.Info(msg)
	}
}

func main() {
	log.Println("Initializing")

	port := os.Getenv("PORT")

	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	g := NewGame()
	g.b.set(3, 4, unit{Name: "hi", Exists: true})

	log.Println("Hello World")
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/game", handlerWrapper(g))
	http.Handle("/ws", websocketWrapper(g))
	log.Info(http.ListenAndServe(":"+port, nil))
}

func handlerWrapper(g Game) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, g.GetStateJSON())
	})
}

func websocketWrapper(g Game) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
		}
		client := Client{conn: conn}
		go client.readPump()
	})
}
