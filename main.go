package main

import (
	"net/http"
	"os"

	"github.com/dchanman/tactics/src/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	log        = logrus.WithField("pkg", "tactics")
	upgrader   = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	mainserver = server.NewServer()
)

func main() {
	log.Info("Initializing")
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/g/{id:[0-9]{6}}", gameHandler)
	router.HandleFunc("/ws/{id:[0-9]{6}}", websocketHandler)

	http.Handle("/", http.FileServer(http.Dir("./webapp/public")))
	http.Handle("/g/", router)
	http.Handle("/ws/", router)
	log.Info(http.ListenAndServe(":"+port, nil))
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.WithFields(logrus.Fields{"id": id}).Info("Websocket request for game")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUpgradeRequired)
		w.Write([]byte("Websocket handshake expected."))
		return
	}
	mainserver.RegisterNewClient(conn)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.WithFields(logrus.Fields{"id": id}).Info("Request for game")
	http.ServeFile(w, r, "./webapp/private/game.html")
}
