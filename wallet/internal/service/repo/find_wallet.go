package repo

import (
	"context"
	db "wallet/internal/database/generated"
	"wallet/internal/service/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *Repo) FindWallet(ctx context.Context, accountID string, hdPath utils.HDPathComponents) (*db.Wallet, error) {
	var result []*db.Wallet
	query := `
		SELECT * 
		FROM wallet
		WHERE account_id = $1 AND
		purpose = $2 AND
		coin_type = $3 AND
		account_index = $4 AND
		change_index = $5 AND
		address_index = $6
		`
	err := pgxscan.Select(ctx, r.db, &result, query,
		accountID, hdPath.Purpose, hdPath.CoinType, hdPath.AccountIndex, hdPath.ChangeIndex, hdPath.AddressIndex)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}
