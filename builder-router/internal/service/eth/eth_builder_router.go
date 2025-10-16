package eth

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service"
	"builder-router/internal/service/validator"
	"context"
	"encoding/json"
	"fmt"
)

type EthBuilderRouter struct {
	builder *EthBuilder
	v       *validator.Validator
}

func NewEthBuilderRouter(builder *EthBuilder) *EthBuilderRouter {
	return &EthBuilderRouter{
		builder: builder,
		v:       validator.NewValidator(),
	}
}

func (e *EthBuilderRouter) WithdrawTx(ctx context.Context, request service.InitiateWithdrawalParams) (txHash string, err error) {
	log := logger.GetLogger("ETH::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return "", fmt.Errorf("invalid params for eth tx: %v", err)
	}
	// unmarshal it to struct
	var ethParams BuildTxEthParams
	err = json.Unmarshal(marshal, &ethParams)
	if err != nil {
		return "", fmt.Errorf("invalid params for eth tx: %v", err)
	}
	hash, err := e.builder.WithdrawTx(ctx, &ethParams)

	if err != nil {
		return "", fmt.Errorf("failed to initiate withdrawal for eth tx: %v", err)
	}
	return hash, nil
}

func (e *EthBuilderRouter) ValidateTx(ctx context.Context, request service.InitiateWithdrawalParams) error {
	log := logger.GetLogger("ETH::WithdrawTx")
	defer func() { log.Debugw("completed") }()
	// parse request params
	details := request.Details
	// marshal them
	marshal, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("invalid params for eth tx: %v", err)
	}
	// unmarshal it to struct
	var ethParams BuildTxEthParams
	err = json.Unmarshal(marshal, &ethParams)
	if err != nil {
		return fmt.Errorf("invalid params for eth tx: %v", err)
	}
	// validate on struct
	err = e.v.Validate(ethParams)
	if err != nil {
		return fmt.Errorf("invalid params for eth tx: %v", err)
	}
	// if erc20 ensure there is contract address
	if ethParams.TxType == "erc20" && ethParams.ContractAddress == "" {
		return fmt.Errorf("invalid params for eth tx, contract address is required for erc20")
	}
	return nil
}
