package main

import (
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

	log.Println("Hello World")
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println(http.ListenAndServe(":"+port, nil))
	log.Println("Hello World Done")
}
