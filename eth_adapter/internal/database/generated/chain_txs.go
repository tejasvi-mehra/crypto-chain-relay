package db

import (
	"time"
)

type ChainTxs struct {
	ID             int        `db:"id"`
	BlockNo        *int       `db:"block_no"`
	BlockTime      *time.Time `db:"block_time"`
	TxType         *string    `db:"tx_type"`
	TxHash         *string    `db:"tx_hash"`
	TraceID        *string    `db:"trace_id"`
	GasLimit       *string    `db:"gas_limit"`
	GasPrice       *string    `db:"gas_price"`
	Amount         *string    `db:"amount"`
	Ticker         *string    `db:"ticker"`
	Spec           *string    `db:"spec"`
	AccountID      *string    `db:"account_id"`
	ActionAddress  *string    `db:"action_address"`
	FromAddress    *string    `db:"from_address"`
	ToAddress      *string    `db:"to_address"`
	AmountInfo     *string    `db:"amount_info"`
	Status         *string    `db:"status"`
	KafkaEventData *string    `db:"kafka_event_data"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}
