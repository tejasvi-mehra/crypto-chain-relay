package btc

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service"
	"builder-router/internal/service/validator"
	"context"
	"encoding/json"
	"fmt"
)

type BtcBuilderRouter struct {
	builder *BtcBuilder
	v       *validator.Validator
}

func NewBtcBuilderRouter(builder *BtcBuilder) *BtcBuilderRouter {
	return &BtcBuilderRouter{
		builder: builder,
		v:       validator.NewValidator(),
	}
}
func (e *BtcBuilderRouter) WithdrawTx(ctx context.Context, request service.InitiateWithdrawalParams) (txHash string, err error) {
	log := logger.GetLogger("BTC::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return "", fmt.Errorf("invalid params for btc tx: %v", err)
	}
	// unmarshal it to struct
	var btcParams BuildTxBtcParams
	err = json.Unmarshal(marshal, &btcParams)
	if err != nil {
		return "", fmt.Errorf("invalid params for btc tx: %v", err)
	}
	hash, err := e.builder.WithdrawTx(ctx, &btcParams)

	if err != nil {
		return "", fmt.Errorf("failed to initiate withdrawal for btc tx: %v", err)
	}
	return hash, nil
}

func (e *BtcBuilderRouter) ValidateTx(ctx context.Context, request service.InitiateWithdrawalParams) error {
	log := logger.GetLogger("BTC::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("invalid params for btc tx: %v", err)
	}
	// unmarshal it to struct
	var btcParams BuildTxBtcParams
	err = json.Unmarshal(marshal, &btcParams)
	if err != nil {
		return fmt.Errorf("invalid params for btc tx: %v", err)
	}
	// validate on struct
	err = e.v.Validate(btcParams)
	if err != nil {
		return fmt.Errorf("invalid params for btc tx: %v", err)
	}
	return nil
}
