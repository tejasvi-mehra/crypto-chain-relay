package transactions

import (
	"context"
	"fmt"
	"transaction-history/internal/framework/logger"
	"transaction-history/internal/service/transactions/repo"
)

type GetAccountBalancesParams struct {
	AccountID string `json:"accountId" validate:"required"`
}

type GetAccountBalancesAggregate struct {
	Amount    float64 `db:"amount" json:"amount"`
	TxType    string  `db:"tx_type" json:"txType"`
	EventType string  `db:"event_type" json:"eventType"`
	Symbol    string  `db:"symbol" json:"symbol"`
}

func (t *TxService) GetAccountBalances(ctx context.Context, params *GetAccountBalancesParams) ([]*repo.AccountBalances, error) {
	log := logger.GetLogger("GetAccountBalances")
	defer func() {
		log.Debugw("completed")
	}()
	if err := t.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return nil, fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)
	balances, err := t.repo.GetAccountBalances(ctx, params.AccountID)
	if err != nil {
		log.Errorw("failed to get transactions for account", "err", err)
		return nil, err
	}
	log.With("count", len(balances))
	if len(balances) == 0 {
		return nil, nil
	}
	return balances, nil
}
