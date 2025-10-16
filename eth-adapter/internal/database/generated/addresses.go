package db

import (
	"time"
)

type Addresses struct {
	ID        int        `db:"id"`
	Address   string     `db:"address"`
	AccountID *string    `db:"account_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
