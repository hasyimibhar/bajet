package main

import (
    "time"

    "github.com/jmoiron/sqlx"
    "github.com/shopspring/decimal"
)

type Item struct {
    Id          int
    Description string
    Cost        decimal.Decimal
    Timestamp   time.Time
}

type ItemList []Item

func (items ItemList) TotalCost() decimal.Decimal {
    total := decimal.NewFromFloat(0)
    for _, i := range items {
        total = total.Add(i.Cost)
    }
    return total
}

type ItemRepository interface {
    All() (ItemList, error)
    Create(item Item) error 
}

type ItemRepo struct {
    db *sqlx.DB
}

func (this *ItemRepo) All() (ItemList, error) {
    items := ItemList{}
    err := this.db.Select(&items, "SELECT * FROM items")
    if err != nil {
        return ItemList{}, err
    }

    return items, nil
}

func (this *ItemRepo) Create(item Item) error {
    _, err := this.db.NamedExec("INSERT INTO items (description, cost) VALUES (:description, :cost)", &item)
    return err
}
