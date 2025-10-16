package fiat

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service"
	"builder-router/internal/service/validator"
	"context"
	"encoding/json"
	"fmt"
)

type FiatBuilderRouter struct {
	builder *FiatBuilder
	v       *validator.Validator
}

func NewFiatBuilderRouter(builder *FiatBuilder) *FiatBuilderRouter {
	return &FiatBuilderRouter{
		builder: builder,
		v:       validator.NewValidator(),
	}
}
func (e *FiatBuilderRouter) WithdrawTx(ctx context.Context, request service.InitiateWithdrawalParams) (txHash string, err error) {
	log := logger.GetLogger("FIAT::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return "", fmt.Errorf("invalid params for fiat tx: %v", err)
	}
	// unmarshal it to struct
	var fiatParams BuildTxFiatParams
	err = json.Unmarshal(marshal, &fiatParams)
	if err != nil {
		return "", fmt.Errorf("invalid params for fiat tx: %v", err)
	}
	hash, err := e.builder.WithdrawTx(ctx, &fiatParams)

	if err != nil {
		return "", fmt.Errorf("failed to initiate withdrawal for fiat tx: %v", err)
	}
	return hash, nil
}

func (e *FiatBuilderRouter) ValidateTx(ctx context.Context, request service.InitiateWithdrawalParams) error {
	log := logger.GetLogger("FIAT::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("invalid params for fiat tx: %v", err)
	}
	// unmarshal it to struct
	var fiatParams BuildTxFiatParams
	err = json.Unmarshal(marshal, &fiatParams)
	if err != nil {
		return fmt.Errorf("invalid params for fiat tx: %v", err)
	}
	// validate on struct
	err = e.v.Validate(fiatParams)
	if err != nil {
		return fmt.Errorf("invalid params for fiat tx: %v", err)
	}
	return nil
}
