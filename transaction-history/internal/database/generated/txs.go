package db

import (
	"time"
)

type Txs struct {
	ID                 int            `db:"id" json:"id"`
	TxType             *string        `db:"tx_type" json:"txType"`
	TraceID            *string        `db:"trace_id" json:"traceId"`
	AmountDecimal      float64        `db:"amount_decimal" json:"amountDecimal"`
	Symbol             *string        `db:"symbol" json:"symbol"`
	AccountID          string         `db:"account_id" json:"accountId"`
	From               *string        `db:"from" json:"from"`
	To                 *string        `db:"to" json:"to"`
	Status             *string        `db:"status" json:"status"`
	TxTime             *time.Time     `db:"tx_time" json:"txTime"`
	EventType          *string        `db:"event_type" json:"eventType"`
	FiatSpecificInfo   map[string]any `db:"fiat_specific_info" json:"fiatSpecificInfo,omitempty"`
	CryptoSpecificInfo map[string]any `db:"crypto_specific_info" json:"cryptoSpecificInfo,omitempty"`
	CreatedAt          *time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt          *time.Time     `db:"updated_at" json:"updatedAt"`
}
