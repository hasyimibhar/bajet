package main

import (
    "os"
    "log"
    "reflect"
    "net/http"
    "html/template"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "github.com/jmoiron/sqlx"
    "github.com/gorilla/schema"
    "github.com/shopspring/decimal"
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

func convertDecimal(value string) reflect.Value {
    number, _ := decimal.NewFromString(value)
    return reflect.ValueOf(number)
}

func AddItemHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    decoder := schema.NewDecoder()
    decoder.RegisterConverter(decimal.NewFromFloat(0), convertDecimal)

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
    databaseUrl := os.Getenv("DATABASE_URL")

    if port == "" {
        log.Fatalln("$PORT must be set")
    }

    db, err := sqlx.Connect("postgres", databaseUrl)
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
