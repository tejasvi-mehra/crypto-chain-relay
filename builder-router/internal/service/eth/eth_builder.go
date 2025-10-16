package eth

import (
	"builder-router/internal/framework/logger"
	"builder-router/internal/service/mocks"
	"builder-router/internal/service/validator"
	"context"
	"fmt"
	"math/big"
	"strings"

	eth_common "github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// this should be a separate service but for convinience sake its here
type EthBuilder struct {
	v       *validator.Validator
	ethNode EthNodeAdapter
	wallet  Wallet
}

type Wallet interface {
	SignTx(ctx context.Context, params *mocks.SignTxParams) (*mocks.SignTxRes, error)
	GetAddressPublicKey(ctx context.Context, params *mocks.GetAddressPublicKeyParams) (*mocks.GetAddressPublicKeyResponse, error)
}

type EthNodeAdapter interface {
	GetNonce() (uint64, error)
	GetGasInfo() (gasLimit uint64, gasPrice *big.Int, err error)
	BroadcastTx(signedTx string) (txHash string, err error)
}

func NewEthBuilder(ethNode EthNodeAdapter, wallet Wallet) *EthBuilder {
	return &EthBuilder{
		v:       validator.NewValidator(),
		ethNode: ethNode,
		wallet:  wallet,
	}
}

type BuildTxEthParams struct {
	AccountId       string  `json:"accountId" validate:"required"`
	ChainId         string  `json:"chainId" validate:"required"`
	TxType          string  `json:"txType" validate:"required"` // "native" or "erc20"
	FromAddress     string  `json:"from" validate:"required"`   // hex address string
	ToAddress       string  `json:"to" validate:"required"`     // hex address string
	Amount          string  `json:"amount" validate:"required"` // string amount in wei (lowest units)
	GasLimit        *uint64 `json:"gasLimit"`
	GasPrice        *string `json:"gasPrice"` // gas price in wei
	Nonce           *uint64 `json:"nonce"`
	ContractAddress string  `json:"contractAddress"` // required if txType = "erc20"
}

func (s *EthBuilder) WithdrawTx(ctx context.Context, params *BuildTxEthParams) (txHash string, err error) {
	log := logger.GetLogger("ETH::WithdrawTx")
	defer func() {
		log.Debugw("completed")
	}()
	// build tx
	unsignedTx, err := s.buildTxEth(params)
	if err != nil {
		return "", err
	}
	// sign tx
	signEthTxParms := &mocks.SignTxParams{
		AccountID:      params.AccountId,
		ChainID:        params.ChainId,
		UnsignedTxData: unsignedTx,
	}
	signedTx, err := s.wallet.SignTx(ctx, signEthTxParms)
	if err != nil {
		return "", err
	}
	// broadcast tx
	txHash, err = s.ethNode.BroadcastTx(signedTx.SignedTx)
	if err != nil {
		return "", err
	}
	log.Infow("broadcasted tx", "hash", txHash)
	return txHash, nil

}

func (s *EthBuilder) buildTxEth(params *BuildTxEthParams) (string, error) {
	log := logger.GetLogger("BuildTxEth")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return "", fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)

	// Parse common fields
	recipient := eth_common.HexToAddress(params.ToAddress)
	var nonce uint64
	if params.Nonce == nil {
		nonce, _ = s.ethNode.GetNonce()
	} else {
		nonce = *params.Nonce
	}

	var gasLimit uint64
	var gasPrice *big.Int
	if params.GasLimit == nil || params.GasPrice == nil {
		gasLimit, gasPrice, _ = s.ethNode.GetGasInfo()
	} else {
		gasLimit = *params.GasLimit
		paramsGasPrice, ok := new(big.Int).SetString(*params.GasPrice, 10)
		if !ok {
			log.Errorw("failed to build tx, invalid gasPrice")
			return "", fmt.Errorf("failed to build tx, invalid gasPrice")
		}
		gasPrice = paramsGasPrice
	}

	var tx *eth_types.Transaction

	if strings.ToLower(params.TxType) == "native" {
		amount, ok := new(big.Int).SetString(params.Amount, 10)
		if !ok {
			log.Errorw("failed to build tx, invalid amount")
			return "", fmt.Errorf("failed to build tx, invalid amount")
		}

		// Native ETH transfer tx
		tx = eth_types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, nil)

	} else if strings.ToLower(params.TxType) == "erc20" {
		if params.ContractAddress == "" {
			log.Errorw("failed to build tx, contractAddress is required")
			return "", fmt.Errorf("failed to build tx, contractAddress is required")
		}
		contractAddress := eth_common.HexToAddress(params.ContractAddress)
		amount, ok := new(big.Int).SetString(params.Amount, 10)
		if !ok {
			log.Errorw("failed to build tx, invalid amount")
			return "", fmt.Errorf("failed to build tx, invalid amount")
		}

		// ERC-20 transfer call data
		methodId := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]
		paddedAddress := eth_common.LeftPadBytes(recipient.Bytes(), 32)
		paddedAmount := eth_common.LeftPadBytes(amount.Bytes(), 32)
		data := append(methodId, append(paddedAddress, paddedAmount...)...)

		tx = eth_types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	} else {
		log.Errorw("failed to build tx, invalid txType")
		return "", fmt.Errorf("failed to build tx, invalid txType")
	}
	unsignedTx, err := tx.MarshalJSON()
	if err != nil {
		log.Errorw("failed to build tx, failed to marshal tx", "err", err)
		return "", fmt.Errorf("failed to build tx, failed to marshal tx")
	}
	return string(unsignedTx), nil

}
