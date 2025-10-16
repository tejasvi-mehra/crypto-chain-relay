package btc

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service/mocks"
	"builder-router/internal/service/validator"
	"context"
	"fmt"
)

// this should be a separate service but for convinience sake its here
type BtcBuilder struct {
	v       *validator.Validator
	btcNode BtcNodeAdapter
	wallet  Wallet
}

type Wallet interface {
	SignTx(ctx context.Context, params *mocks.SignTxParams) (*mocks.SignTxRes, error)
	GetAddressPublicKey(ctx context.Context, params *mocks.GetAddressPublicKeyParams) (*mocks.GetAddressPublicKeyResponse, error)
}

type BtcNodeAdapter interface {
	GetNonce() (uint64, error)
	GetFeeInfo() (feeRete string, err error)
	BroadcastTx(signedTx string) (txHash string, err error)
}

func NewBtcBuilder(btcNode BtcNodeAdapter, wallet Wallet) *BtcBuilder {
	return &BtcBuilder{
		v:       validator.NewValidator(),
		btcNode: btcNode,
		wallet:  wallet,
	}
}

type BuildTxBtcParams struct {
	AccountId     string  `json:"accountId" validate:"required"`
	ChainId       string  `json:"chainId" validate:"required"`
	From          string  `json:"from" validate:"required"`
	To            string  `json:"to" validate:"required"`
	Amount        string  `json:"amount" validate:"required"`
	FeeRate       *string `json:"feeRate,omitempty"`
	Nonce         *uint64 `json:"nonce,omitempty"`
	ChangeAddress *string `json:"changeAddress,omitempty"`
}

func (s *BtcBuilder) WithdrawTx(ctx context.Context, params *BuildTxBtcParams) (txHash string, err error) {
	log := logger.GetLogger("BTC::WithdrawTx")
	defer func() {
		log.Debugw("completed")
	}()
	// build tx
	psbt, err := s.buildTxBtc(params)
	if err != nil {
		return "", err
	}
	// sign tx
	signBtcTxParms := &mocks.SignTxParams{
		AccountID:      params.AccountId,
		ChainID:        params.ChainId,
		UnsignedTxData: psbt,
	}
	signedTx, err := s.wallet.SignTx(ctx, signBtcTxParms)
	if err != nil {
		return "", err
	}
	// broadcast tx
	txHash, err = s.btcNode.BroadcastTx(signedTx.SignedTx)
	if err != nil {
		return "", err
	}
	log.Infow("broadcasted tx", "hash", txHash)
	return txHash, nil

}

func (s *BtcBuilder) buildTxBtc(params *BuildTxBtcParams) (string, error) {
	log := logger.GetLogger("BuildTxBtc")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return "", fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)

	// check for nonce or fetch from node
	// check if UTXOs can be parsed in correct way
	// check if feeRate is empty then fetch from node
	// check feeRate is more than dust
	// return psbt
	return "params-psbt", nil
}
