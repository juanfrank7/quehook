package main

import (
	"log"
	"net/http"

	"google.golang.org/appengine"

	"github.com/forstmeier/watchmyrepo/database"
	"github.com/forstmeier/watchmyrepo/handlers"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handlers.Display("site/", db))
	http.HandleFunc("/submit", handlers.Submit("site/", db))
	http.HandleFunc("/data", handlers.Data(db))
	http.HandleFunc("/faq", handlers.FAQ("site/"))
	http.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/main.css")
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "site/favicon.ico")
	})

	appengine.Main()
}
