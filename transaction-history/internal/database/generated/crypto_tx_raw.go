package db

import (
	"time"
)

type CryptoTxRaw struct {
	ID        int        `db:"id"`
	TraceID   *string    `db:"trace_id"`
	EventData *string    `db:"event_data"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
