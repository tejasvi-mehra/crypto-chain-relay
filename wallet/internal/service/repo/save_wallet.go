package repo

import (
	"context"
	"time"
	db "wallet/internal/database/generated"
	"wallet/internal/framework/logger"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *Repo) SaveWallet(ctx context.Context, wallet db.Wallet) error {
	log := logger.GetLogger("SaveWallet")
	defer log.Sync()

	var id []string
	if err := pgxscan.Select(ctx, r.db, &id,
		`INSERT INTO wallet (
			account_id,
			seed_id,
			purpose,
			coin_type,
			account_index,
			change_index,
			address_index,
			address,
			public_key,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (seed_id, purpose, coin_type, account_index, change_index, address_index)
		DO NOTHING`,
		wallet.AccountID, wallet.SeedID, wallet.Purpose, wallet.CoinType, wallet.AccountIndex, wallet.ChangeIndex, wallet.AddressIndex, wallet.Address, wallet.PublicKey, time.Now().UTC(), time.Now().UTC(),
	); err != nil {
		log.Errorw("fail to save enterprise info", "error", err)
		return err
	}
	return nil
}
