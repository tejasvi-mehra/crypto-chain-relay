package service

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service/mocks"
	"context"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// RequestBody defines the generic structure for withdrawal requests
type InitiateWithdrawalParams struct {
	WithdrawalType string         `json:"withdrawalType" validate:"required,oneof=eth btc fiat"`
	AccountId      string         `json:"accountId" validate:"required"`
	Amount         string         `json:"amount" validate:"required"`
	Symbol         string         `json:"symbol" validate:"required"`
	Details        map[string]any `json:"details" validate:"required"`
}

type InitiateWithdrawalResponse struct {
	Type                string `json:"type"`
	WithdrawalReference string `json:"withdrawalReference"`
}

func (s *Service) InitiateWithdrawal(ctx context.Context, e echo.Context, params *InitiateWithdrawalParams) (*InitiateWithdrawalResponse, error) {
	log := logger.GetLogger("GetAddressPublicKey")
	defer func() {
		log.Debugw("completed")
	}()
	var err error

	err = s.ValidateWithdrawal(ctx, e, params)
	if err != nil {
		return nil, err
	}
	log.With("params", params)
	// do additional validations
	// interface call
	var txReference string
	switch params.WithdrawalType {
	case "eth":
		txReference, err = s.eth.WithdrawTx(ctx, *params)
	case "btc":
		txReference, err = s.btc.WithdrawTx(ctx, *params)
	case "fiat":
		txReference, err = s.btc.WithdrawTx(ctx, *params)
	default:
		return nil, fmt.Errorf("invalid withdrawal type: %s", params.WithdrawalType)
	}
	if err != nil {
		log.Errorw("failed to initiate withdrawal", "err", err)
		return nil, err
	}
	return &InitiateWithdrawalResponse{
		WithdrawalReference: txReference,
		Type:                params.WithdrawalType,
	}, nil

}

func (s *Service) ValidateWithdrawal(ctx context.Context, e echo.Context, params *InitiateWithdrawalParams) error {
	log := logger.GetLogger("GetAddressPublicKey")
	defer func() {
		log.Debugw("completed")
	}()
	// params validation
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return fmt.Errorf("invalid params: %v", err)
	}
	// get JWT from context
	user := e.Get("user")
	token, ok := user.(*jwt.Token)
	if !ok {
		return fmt.Errorf("invalid token")
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return fmt.Errorf("no subject in token")
	}
	if subject != params.AccountId {
		return fmt.Errorf("unauthorized, param accountID doesn't match JWT")
	}
	// parse amount to see greater than 0
	amount, err := strconv.ParseFloat(params.Amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount, must be greater than 0: %v", err)
	}
	// do balance check for symbol used
	getBalanceParams := &mocks.GetAccountBalancesParams{
		AccountID: params.AccountId,
	}
	balances, err := s.txHistory.GetAccountBalances(ctx, getBalanceParams)
	if err != nil {
		return fmt.Errorf("failed to get balance for withdrawal: %v", err)
	}
	if len(balances) == 0 {
		return fmt.Errorf("cannot init withdrawal, no balance for account")
	}

	var balance float64 = 0
	for _, b := range balances {
		// calculate balance for symbol
		if b.Symbol == params.Symbol {
			if b.TxType == "withdrawal" || b.TxType == "network_fee" {
				balance -= b.Amount
			} else {
				balance += b.Amount
			}
		}
	}
	// found balance for symbol
	if balance <= amount {
		return fmt.Errorf("cannot init withdrawal, not enough balance: %s, available: %f, withdrawalAmount: %f", params.Symbol, balance, amount)
	} else {
		// has balance continue
		return nil
	}
}
