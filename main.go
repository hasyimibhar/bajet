package main

import (
    "os"
    "log"
    "fmt"
    "net/http"
    "html/template"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    "github.com/gorilla/schema"
)

var Items ItemRepository

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    items, err := Items.All()
    if err != nil {
        log.Println("[ERROR]", err.Error())
        w.Write([]byte(err.Error()))
        return
    }

    t, _ := template.ParseFiles("public/index.html")
    t.Execute(w, items)
}

func AddItemHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    decoder := schema.NewDecoder()
    item := Item{}
    decoder.Decode(&item, r.PostForm)

    err := Items.Create(item)
    if err != nil {
        log.Println("[ERROR]", err.Error())
        w.Write([]byte(err.Error()))
        return
    }

    http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
    port := os.Getenv("PORT")
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    dbPassword := os.Getenv("DB_PASSWORD")

    if port == "" {
        log.Fatalln("$PORT must be set")
    }

    db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable",
        dbUser, dbName, dbPassword))
    if err != nil {
        log.Fatalln(err)
    }

    defer db.Close()
    Items = &ItemRepo{db}

    r := mux.NewRouter()
    r.HandleFunc("/", IndexHandler).Methods(http.MethodGet)
    r.HandleFunc("/items", AddItemHandler).Methods(http.MethodPost)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.ListenAndServe(":" + port, r)
}
