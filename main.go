package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"

	"github.com/forstmeier/watchmyrepo/database"
	"github.com/forstmeier/watchmyrepo/handlers"
)

type server struct {
	secret string
	server http.Server
	client *github.Client
}

func newServer(secret, addr string, c *github.Client, database database.DB) server {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Display("site/", database)).Methods("GET") //.Schemes("https")
	r.HandleFunc("/submit", handlers.Submit("site/", database)).Methods("POST").Schemes("https")
	r.HandleFunc("/data", handlers.Data(database)).Methods("POST").Schemes("https")
	r.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/main.css")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/favicon.ico")
	})

	s := server{
		secret: secret,
		server: http.Server{
			Addr:         addr,
			Handler:      r,
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		},
		client: c,
	}

	return s
}

func (s *server) start() {
	s.server.ListenAndServe()
}

func (s *server) stop() {
	s.server.Close()
}

func main() {
	boltDB, err := bolt.Open("watchmyrepo.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer boltDB.Close()

	db, err := database.New(boltDB)
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient(nil)
	port := os.Getenv("PORT")
	if port == "" {
		port = "2187"
	}
	fmt.Println("PORT:", port)

	server := newServer("secret", ":"+port, client, db)
	server.start()
}
