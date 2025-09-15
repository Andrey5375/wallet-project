package models

type Wallet struct {
	ID      string  `db:"id" json:"id"`
	Balance float64 `db:"balance" json:"balance"`
}
