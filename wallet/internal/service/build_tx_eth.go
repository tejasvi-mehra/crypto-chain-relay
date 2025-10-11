package service

// this is just a sample code for easy generation of eth tx in a way accepted by the signTx endpoint - this should be on a separate service tx-builder-eth. tx-builder should route eth related txs to tx-builder-eth to the same endpoint

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"wallet/internal/framework/logger"

	eth_common "github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type BuildTxEthParams struct {
	ChainId         string `json:"chainId" validate:"required"`
	TxType          string `json:"txType" validate:"required"`    // "native" or "erc20"
	ToAddress       string `json:"toAddress" validate:"required"` // hex address string
	Amount          string `json:"amount" validate:"required"`    // string amount in wei (lowest units)
	GasLimit        uint64 `json:"gasLimit" validate:"required"`
	GasPrice        string `json:"gasPrice" validate:"required"` // gas price in wei
	Nonce           uint64 `json:"nonce" validate:"required"`
	ContractAddress string `json:"contractAddress"` // required if txType = "erc20"
}

type BuildTxEthResponse struct {
	UnsignedTx string `json:"unsignedTx"` // marshalled unsigned tx
}

func (s *Service) BuildTxEth(ctx context.Context, params *BuildTxEthParams) (*BuildTxEthResponse, error) {
	log := logger.GetLogger("BuildTxEth")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return nil, fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)

	// Parse common fields
	recipient := eth_common.HexToAddress(params.ToAddress)
	gasLimit := params.GasLimit
	gasPrice, ok := new(big.Int).SetString(params.GasPrice, 10)
	if !ok {
		log.Errorw("failed to build tx, invalid gasPrice")
		return nil, fmt.Errorf("failed to build tx, invalid gasPrice")
	}
	nonce := params.Nonce

	var tx *eth_types.Transaction

	if strings.ToLower(params.TxType) == "native" {
		amount, ok := new(big.Int).SetString(params.Amount, 10)
		if !ok {
			log.Errorw("failed to build tx, invalid amount")
			return nil, fmt.Errorf("failed to build tx, invalid amount")
		}

		// Native ETH transfer tx
		tx = eth_types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, nil)

	} else if strings.ToLower(params.TxType) == "erc20" {
		if params.ContractAddress == "" {
			log.Errorw("failed to build tx, contractAddress is required")
			return nil, fmt.Errorf("failed to build tx, contractAddress is required")
		}
		contractAddress := eth_common.HexToAddress(params.ContractAddress)
		amount, ok := new(big.Int).SetString(params.Amount, 10)
		if !ok {
			log.Errorw("failed to build tx, invalid amount")
			return nil, fmt.Errorf("failed to build tx, invalid amount")
		}

		// ERC-20 transfer call data
		methodId := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]
		paddedAddress := eth_common.LeftPadBytes(recipient.Bytes(), 32)
		paddedAmount := eth_common.LeftPadBytes(amount.Bytes(), 32)
		data := append(methodId, append(paddedAddress, paddedAmount...)...)

		tx = eth_types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	} else {
		log.Errorw("failed to build tx, invalid txType")
		return nil, fmt.Errorf("failed to build tx, invalid txType")
	}
	unsignedTx, err := tx.MarshalJSON()
	if err != nil {
		log.Errorw("failed to build tx, failed to marshal tx", "err", err)
		return nil, fmt.Errorf("failed to build tx, failed to marshal tx")
	}
	return &BuildTxEthResponse{
		UnsignedTx: string(unsignedTx),
	}, nil
}
