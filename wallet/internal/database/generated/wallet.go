package db

import (
	"time"
)

type Wallet struct {
	ID           int       `db:"id"`
	AccountID    string    `db:"account_id"`
	SeedID       string    `db:"seed_id"`
	Purpose      int       `db:"purpose"`
	CoinType     int       `db:"coin_type"`
	AccountIndex int       `db:"account_index"`
	ChangeIndex  int       `db:"change_index"`
	AddressIndex int       `db:"address_index"`
	PublicKey    string    `db:"public_key"`
	Address      string    `db:"address"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
