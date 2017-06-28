package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("pkg", "tactics")
)

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
	log.Info(http.ListenAndServe(":"+port, nil))
}

func handlerWrapper(g Game) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, g.GetStateJSON())
	})
}
