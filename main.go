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
	router.HandleFunc("/g/", gameHandler)
	router.HandleFunc("/g/{id:[0-9][0-9][0-9][0-9][0-9][0-9]}", gameHandler)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocketWrapper())
	http.Handle("/g/", router)
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

func gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Write([]byte("Hello: "))
	w.Write([]byte(string(id)))
}
