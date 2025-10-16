package db

import (
	"time"
)

type BlockData struct {
	ID             int       `db:"id"`
	BlockNo        int64     `db:"block_no"`
	BlockHash      string    `db:"block_hash"`
	ParentHash     string    `db:"parent_hash"`
	BlockTimestamp time.Time `db:"block_timestamp"`
	ProcessTimeMs  int64     `db:"process_time_ms"`
	Confirmed      bool      `db:"confirmed"`
	RetryCount     int       `db:"retry_count"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
