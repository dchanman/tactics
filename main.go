package main

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	log        = logrus.WithField("pkg", "tactics")
	upgrader   = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	mainserver = NewServer()
)

func main() {
	log.Println("Initializing")

	port := os.Getenv("PORT")

	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocketWrapper())
	log.Info(http.ListenAndServe(":"+port, nil))
}

func websocketWrapper() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusUpgradeRequired)
			w.Write([]byte("Websocket handshake expected."))
			return
		}
		mainserver.RegisterNewClient(conn)
	})
}
