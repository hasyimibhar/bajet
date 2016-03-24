package main

import (
    "net/http"

    "github.com/gorilla/mux"
)

func AddItemHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/items", AddItemHandler).Methods(http.MethodPost)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.ListenAndServe(":8000", r)
}
