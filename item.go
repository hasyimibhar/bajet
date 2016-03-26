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

func (item Item) CostString() string {
    return item.Cost.StringFixed(2)
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
    Update(id int, data map[string]interface{}) error
}

type ItemRepo struct {
    db *sqlx.DB
}

func (this *ItemRepo) All() (ItemList, error) {
    items := ItemList{}
    err := this.db.Select(&items, "SELECT * FROM items ORDER BY id ASC")
    if err != nil {
        return ItemList{}, err
    }

    return items, nil
}

func (this *ItemRepo) Create(item Item) error {
    _, err := this.db.NamedExec("INSERT INTO items (description, cost) VALUES (:description, :cost)", &item)
    return err
}

func (this *ItemRepo) Update(id int, data map[string]interface{}) error {
    data["id"] = id
    _, err := this.db.NamedExec("UPDATE items SET description=:description,cost=:cost WHERE id=:id", data)
    return err
}
