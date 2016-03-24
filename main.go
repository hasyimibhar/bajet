package main

import (
    "log"
    "net/http"
    "html/template"

    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
)

type Item struct {
    Description string
    Cost        int
}

var Items []Item

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("public/index.html")
    t.Execute(w, Items)
}

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
    r.HandleFunc("/", IndexHandler).Methods(http.MethodGet)
    r.HandleFunc("/items", AddItemHandler).Methods(http.MethodPost)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.ListenAndServe(":8000", r)
}
