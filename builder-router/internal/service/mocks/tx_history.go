package mocks

import "context"

type GetAccountBalancesParams struct {
	AccountID string `json:"accountId" validate:"required"`
}

type GetAccountBalancesRes struct {
	Amount    float64 `json:"amount"`
	TxType    string  `json:"txType"`
	EventType string  `json:"eventType"`
	Symbol    string  `json:"symbol"`
}

type TransactionHistory struct {
}

func (t *TransactionHistory) GetAccountBalances(ctx context.Context, params *GetAccountBalancesParams) ([]*GetAccountBalancesRes, error) {
	balances := []*GetAccountBalancesRes{
		{
			Amount:    10,
			TxType:    "withdrawal",
			EventType: "crypto",
			Symbol:    "ETH",
		},
		{
			Amount:    20,
			TxType:    "deposit",
			EventType: "crypto",
			Symbol:    "ETH",
		},
		{
			Amount:    10,
			TxType:    "withdrawal",
			EventType: "crypto",
			Symbol:    "BTC",
		},
		{
			Amount:    20,
			TxType:    "deposit",
			EventType: "crypto",
			Symbol:    "BTC",
		},
	}
	return balances, nil
}
