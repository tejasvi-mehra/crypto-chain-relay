package service

import (
	"builder-router/internal/service/mocks"
	"builder-router/internal/service/validator"
	"context"
)

type Service struct {
	eth       EthAdapter
	btc       BTCAdapter
	fiat      FiatAdapter
	txHistory TransactionHistory
	v         *validator.Validator
}

type TransactionHistory interface {
	GetAccountBalances(ctx context.Context, params *mocks.GetAccountBalancesParams) ([]*mocks.GetAccountBalancesRes, error)
}

type EthAdapter interface {
	WithdrawTx(ctx context.Context, request InitiateWithdrawalParams) (txHash string, err error)
}
type BTCAdapter interface {
	WithdrawTx(ctx context.Context, request InitiateWithdrawalParams) (txHash string, err error)
}
type FiatAdapter interface {
	WithdrawTx(ctx context.Context, request InitiateWithdrawalParams) (referenceID string, err error)
}

func NewService(eth EthAdapter, btc BTCAdapter, fiat FiatAdapter, txHistory TransactionHistory) *Service {
	return &Service{
		v:         validator.NewValidator(),
		eth:       eth,
		btc:       btc,
		fiat:      fiat,
		txHistory: txHistory,
	}

}
