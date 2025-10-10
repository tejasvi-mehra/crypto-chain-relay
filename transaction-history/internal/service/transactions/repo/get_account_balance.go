package repo

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type AccountBalances struct {
	Amount    float64 `db:"amount" json:"amount"`
	TxType    string  `db:"tx_type" json:"txType"`
	EventType string  `db:"event_type" json:"eventType"`
	Symbol    string  `db:"symbol" json:"symbol"`
}

// Get transactions for an account by accountID with pagination
func (r *Repo) GetAccountBalances(ctx context.Context, accountID string) ([]*AccountBalances, error) {
	var result []*AccountBalances
	query := `
		SELECT SUM(amount_decimal) as amount,tx_type, event_type, symbol FROM txs WHERE account_id = $1 GROUP BY tx_type, event_type, symbol`
	err := pgxscan.Select(ctx, r.db, &result, query,
		accountID)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}
