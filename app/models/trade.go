package models

import "time"

type Trade struct {
	ID          int
	BuyOrderID  int
	Buyer       string
	SellOrderID int
	Seller      string
	Amount      float64
	Price       float64
	CreatedAt   time.Time
}
