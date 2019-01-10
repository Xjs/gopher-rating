package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Xjs/gopher-rating/storage/mysql"
	"github.com/Xjs/gopher-rating/view"
	"github.com/gorilla/mux"
)

func main() {
	storer, err := mysql.NewStorer(context.Background(), os.Getenv("GOPHER_STORER_MYSQL"))
	if err != nil {
		log.Fatal(err)
	}
	var template *template.Template

	handler := view.NewHandler(storer, template)

	listen := os.Getenv("GOPHER_STORER_LISTEN")
	if listen == "" {
		listen = ":8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Gophers)
	r.HandleFunc("/gopher/{hash:[0-9a-f]+}", handler.Gopher)
	r.HandleFunc("/rate/{hash:[0-9a-f]+}/{rating:[0-9]+}", handler.Rate)
	r.HandleFunc("/upload", handler.Upload).Methods("POST", "PUT")
	http.Handle("/", r)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatal(err)
	}
}
