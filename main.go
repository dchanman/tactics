package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dchanman/tactics/src/server"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	log        = logrus.WithField("pkg", "tactics")
	upgrader   = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	mainserver = server.NewServer()

	templates = template.Must(template.ParseFiles(
		"./webapp/private/game.tmpl",
		"./webapp/private/howtoplay.tmpl",
		"./webapp/private/lobby.tmpl",
		"./webapp/private/partials/chat.tmpl",
		"./webapp/private/partials/commonjs.tmpl",
		"./webapp/private/partials/head.tmpl",
		"./webapp/private/partials/history.tmpl",
		"./webapp/private/partials/nav.tmpl",
		"./webapp/private/partials/overlaysetting.tmpl",
		"./webapp/private/partials/status.tmpl"))
)

func main() {
	log.Info("Initializing")

	rand.Seed(time.Now().UTC().UnixNano())
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	// Game ID routes
	router := mux.NewRouter()
	router.HandleFunc("/g/{id:[0-9]{6}}", gameHandler)

	// JSON-RPC hooks
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterService(mainserver, "")

	http.Handle("/g/", router)
	http.HandleFunc("/ws", websocketHandler)
	http.Handle("/data/", rpcServer)
	http.Handle("/", blockDirListing(http.FileServer(http.Dir("./webapp/public"))))

	log.Info(http.ListenAndServe(":"+port, nil))
}

func blockDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templates.ExecuteTemplate(w, "lobby", nil)
			return
		}
		if r.URL.Path == "/howtoplay" {
			templates.ExecuteTemplate(w, "howtoplay", nil)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getMuxGameID(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	gameid, err := strconv.ParseUint(id, 10, 32)
	return uint32(gameid), err
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
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
	gameid, err := getMuxGameID(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if !mainserver.DoesGameIDExist(gameid) {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "game", nil)
}
