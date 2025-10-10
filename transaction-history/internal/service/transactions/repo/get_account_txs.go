package repo

import (
	"context"

	db "transaction-history/internal/database/generated"

	"github.com/georgysavva/scany/v2/pgxscan"
)

// Get transactions for an account by accountID with pagination
func (r *Repo) GetAccountTransactions(ctx context.Context, accountID string, limit, offset int) ([]*db.Txs, error) {
	var result []*db.Txs
	query := `
		SELECT * 
		FROM txs
		WHERE account_id = $1
		ORDER BY tx_time DESC
		LIMIT $2 OFFSET $3`
	err := pgxscan.Select(ctx, r.db, &result, query,
		accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}
