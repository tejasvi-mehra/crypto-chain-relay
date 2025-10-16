package fiat

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service/validator"
	"context"
	"fmt"
)

// this should be a separate service but for convinience sake its here
type FiatBuilder struct {
	v *validator.Validator
}

func NewFiatBuilder() *FiatBuilder {
	return &FiatBuilder{
		v: validator.NewValidator(),
	}
}

// FIATWithdrawalDetails holds details specific to FIAT withdrawals
type BuildTxFiatParams struct {
	BankAccountNumber string  `json:"bankAccountNumber" validate:"required"`
	BankName          string  `json:"bankName" validate:"required"`
	Currency          string  `json:"currency" validate:"required"`
	Reference         *string `json:"reference,omitempty"`
}

func (s *FiatBuilder) WithdrawTx(ctx context.Context, params *BuildTxFiatParams) (referenceID string, err error) {
	log := logger.GetLogger("FIAT::WithdrawTx")
	defer func() {
		log.Debugw("completed")
	}()
	// complete fiat withdrawal
	referenceID, err = s.initateWithdrawal(params)
	if err != nil {
		return "", err
	}

	return referenceID, nil

}

func (s *FiatBuilder) initateWithdrawal(params *BuildTxFiatParams) (string, error) {
	log := logger.GetLogger("BuildTxFiat")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return "", fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)
	// completed withdrawal
	return "completed-witdrawal", nil
}
