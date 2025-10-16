package builder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ETHWithdrawalDetails holds details specific to ETH withdrawals
type ETHWithdrawalDetails struct {
	From            string  `json:"from" validate:"required"`
	ChainId         string  `json:"chainId" validate:"required"`
	TxType          string  `json:"txType" validate:"required,oneof=native erc20"`
	To              string  `json:"to" validate:"required"`
	Amount          string  `json:"amount" validate:"required"`
	GasLimit        *uint64 `json:"gasLimit,omitempty"`
	GasPrice        *string `json:"gasPrice,omitempty"`
	Nonce           *uint64 `json:"nonce,omitempty"`
	ContractAddress *string `json:"contractAddress,omitempty"`
}

// BTCWithdrawalDetails holds details specific to BTC withdrawals
type BTCWithdrawalDetails struct {
	ChainId       string  `json:"chainId" validate:"required"`
	From          string  `json:"from" validate:"required"`
	To            string  `json:"to" validate:"required"`
	Amount        string  `json:"amount" validate:"required"`
	FeeRate       *string `json:"feeRate,omitempty"`
	Nonce         *uint64 `json:"nonce,omitempty"`
	ChangeAddress *string `json:"changeAddress,omitempty"`
}

// FIATWithdrawalDetails holds details specific to FIAT withdrawals
type FIATWithdrawalDetails struct {
	BankAccountNumber string  `json:"bankAccountNumber" validate:"required"`
	BankName          string  `json:"bankName" validate:"required"`
	Currency          string  `json:"currency" validate:"required"`
	Reference         *string `json:"reference,omitempty"`
}

// Builder holds the validator and builds withdrawal requests
type Builder struct {
	validator echo.Validator
}

func NewBuilder(v echo.Validator) *Builder {
	return &Builder{validator: v}
}

// RequestBody defines the generic structure for withdrawal requests
type RequestBody struct {
	WithdrawalType string          `json:"withdrawalType" validate:"required,oneof=eth btc fiat"`
	AccountId      string          `json:"accountId" validate:"required"`
	Amount         string          `json:"amount" validate:"required"`
	Details        json.RawMessage `json:"details" validate:"required"`
}

// Common build function dispatching based on withdrawal type
func (b *Builder) Build(req RequestBody) (map[string]string, error) {
	if err := b.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("failed to build tx: %v", err)
	}
	switch req.WithdrawalType {
	case "eth":
		var ethDetails ETHWithdrawalDetails
		if err := json.Unmarshal(req.Details, &ethDetails); err != nil {
			return nil, err
		}
		return b.BuildETH(req.AccountId, req.Amount, ethDetails)

	case "btc":
		var btcDetails BTCWithdrawalDetails
		if err := json.Unmarshal(req.Details, &btcDetails); err != nil {
			return nil, err
		}
		return b.BuildBTC(req.AccountId, req.Amount, btcDetails)

	case "fiat":
		var fiatDetails FIATWithdrawalDetails
		if err := json.Unmarshal(req.Details, &fiatDetails); err != nil {
			return nil, err
		}
		return b.BuildFIAT(req.AccountId, req.Amount, fiatDetails)

	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "unsupported withdrawal type")
	}
}
func (b *Builder) BuildETH(accountId, amount string, details ETHWithdrawalDetails) (map[string]string, error) {
	if err := b.validator.Validate(details); err != nil {
		return nil, err
	}
	// TODO: Implement actual ETH withdrawal logic here
	return map[string]string{"transactionId": "eth123", "status": "accepted"}, nil
}

func (b *Builder) BuildBTC(accountId, amount string, details BTCWithdrawalDetails) (map[string]string, error) {
	if err := b.validator.Validate(details); err != nil {
		return nil, err
	}
	// TODO: Implement actual BTC withdrawal logic here
	return map[string]string{"transactionId": "btc123", "status": "accepted"}, nil
}

func (b *Builder) BuildFIAT(accountId, amount string, details FIATWithdrawalDetails) (map[string]string, error) {
	if err := b.validator.Validate(details); err != nil {
		return nil, err
	}
	// TODO: Implement actual FIAT withdrawal logic here
	return map[string]string{"transactionId": "fiat123", "status": "accepted"}, nil
}
