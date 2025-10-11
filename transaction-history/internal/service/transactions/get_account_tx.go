package transactions

import (
	"context"
	"fmt"
	db "transaction-history/internal/database/generated"
	"transaction-history/internal/framework/logger"
)

type GetAccountTransactionsParams struct {
	AccountID string `json:"accountId" validate:"required"`
	Limit     int    `json:"limit" validate:"gt=0"`
	Offset    int    `json:"offset"`
}

func (t *TxService) GetAccountTransactions(ctx context.Context, params *GetAccountTransactionsParams) ([]*db.Txs, error) {
	log := logger.GetLogger("GetAccountTransactions")
	defer func() {
		log.Debugw("completed")
	}()
	if err := t.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return nil, fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)
	txs, err := t.repo.GetAccountTransactions(ctx, params.AccountID, params.Limit, params.Offset)
	if err != nil {
		log.Errorw("failed to get transactions for account", "err", err)
		return nil, err
	}
	log.With("count", len(txs))
	return txs, nil
}
