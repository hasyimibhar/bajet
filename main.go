package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
)

type Item struct {
    Description string
    Cost        int
}

var Items []Item

func AddItemHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    decoder := schema.NewDecoder()
    item := Item{}
    decoder.Decode(&item, r.PostForm)

    Items = append(Items, []Item{item}...)
    log.Println("[TRACE] Items:", Items)

    http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
    Items = []Item{}

    r := mux.NewRouter()
    r.HandleFunc("/items", AddItemHandler).Methods(http.MethodPost)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.ListenAndServe(":8000", r)
}
