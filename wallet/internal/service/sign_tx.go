package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"wallet/internal/framework/logger"

	"wallet/internal/service/utils"

	"github.com/btcsuite/btcd/btcutil/psbt"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type SignTxParams struct {
	AccountID      string `json:"accountId" validate:"required"`
	ChainID        string `json:"chainId" validate:"required"`
	UnsignedTxData any    `json:"unsignedTxData" validate:"required"`
}

type UnsignedTxDataEth struct {
	UnsignedTx string `json:"unsignedTx"`
}

type UnsignedTxDataBTC struct {
	Psbt string `json:"psbt"`
}

type SignTxResponse struct {
	SignedTx string `json:"signedTx"` // hex-encoded signed tx
}

func (s *Service) SignTx(ctx context.Context, params *SignTxParams) (*SignTxResponse, error) {
	log := logger.GetLogger("SignTx")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return nil, fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)

	hdPath, ok := utils.HDPaths[params.ChainID]
	if !ok {
		log.Errorw("invalid chainID, no HDPath")
		return nil, fmt.Errorf("invalid chainID, no hdPath")
	}
	chainType, ok := utils.ChainMap[params.ChainID]
	if !ok {
		log.Errorw("invalid chainID, unknown chainType")
		return nil, fmt.Errorf("invalid chainID, unknown chainType")
	}
	wallet, err := s.repo.FindWallet(ctx, params.AccountID, hdPath)
	if err != nil {
		// this error does not include no row found, need to explicitly check for that
		log.Errorw("failed to find wallet", "err", err)
		return nil, fmt.Errorf("failed to find wallet for accountID")
	}
	if wallet == nil {
		log.Errorw("no wallet found for given accountID and chainID", "hdPath", hdPath, "chainID", params.ChainID)
		return nil, fmt.Errorf("no wallet found for given accountID and chainID")
	}

	seedUUID, err := uuid.Parse(wallet.SeedID)
	if err != nil {
		// this error does not include no row found, need to explicitly check for that
		log.Errorw("invalid seedID for wallet", "err", err)
		return nil, fmt.Errorf("failed to find wallet for accountID, invalid seedID")
	}
	log.Debugw("chainType", "chainType", chainType)
	var signedTx string
	if chainType == utils.ETH {
		var ethTxData UnsignedTxDataEth
		// parse unsigned chain type data as eth type data
		marshal, err := json.Marshal(params.UnsignedTxData)
		if err != nil {
			log.Errorw("invalid unsigned tx data for eth tx, cannot marshal", "type", reflect.TypeOf(params.UnsignedTxData))
			return nil, fmt.Errorf("invalid unsigned tx data for eth tx,cannot marshal")
		}
		err = json.Unmarshal(marshal, &ethTxData)
		if err != nil {
			log.Errorw("invalid unsigned tx data for eth tx, cannot unmarshal", "type", reflect.TypeOf(params.UnsignedTxData))
			return nil, fmt.Errorf("invalid unsigned tx data for eth tx,cannot unmarshal")
		}
		log.Debugw("got txData", "data", ethTxData)
		// validate
		if ethTxData.UnsignedTx == "" {
			log.Errorw("invalid unsigned tx data for eth tx, empty unsigned tx")
			return nil, fmt.Errorf("invalid unsigned tx data for eth tx, empty unsigned tx")
		}
		// parse as ETH tx
		var txEth eth_types.Transaction
		err = txEth.UnmarshalJSON([]byte(ethTxData.UnsignedTx))
		if err != nil {
			log.Errorw("invalid unsigned tx data for eth tx, failed to unmarshal", "err", err)
			return nil, fmt.Errorf("invalid unsigned tx data for eth tx, failed to unmarshal")
		}
		log.Debugw("got eth tx", "txEth", &txEth)
		ethSignTx, err := s.kms.SignTx(seedUUID, params.ChainID, hdPath, &txEth)
		if err != nil {
			// this error does not include no row found, need to explicitly check for that
			log.Errorw("failed to sign tx", "err", err)
			return nil, fmt.Errorf("failed to sign tx")
		}
		signedTx = "0x" + hex.EncodeToString(ethSignTx)

	} else if chainType == utils.BTC {
		log.With("chainType", utils.BTC)
		var btcTxData UnsignedTxDataBTC
		// parse unsigned chain type data as btc type data
		marshal, err := json.Marshal(params.UnsignedTxData)
		if err != nil {
			log.Errorw("invalid unsigned tx data for btc tx, cannot marshal", "type", reflect.TypeOf(params.UnsignedTxData))
			return nil, fmt.Errorf("invalid unsigned tx data for btc tx,cannot marshal")
		}
		log.Debugw("marshalled")
		err = json.Unmarshal(marshal, &btcTxData)
		if err != nil {
			log.Errorw("invalid unsigned tx data for btc tx, cannot unmarshal", "type", reflect.TypeOf(params.UnsignedTxData))
			return nil, fmt.Errorf("invalid unsigned tx data for btc tx,cannot unmarshal")
		}
		log.Debugw("got txData", "data", btcTxData)
		btcPsbt, err := psbt.NewFromRawBytes(
			strings.NewReader(btcTxData.Psbt), true,
		)
		if err != nil {
			log.Errorw("invalid unsigned tx data for btc tx, cannot get psbt", "err", err)
			return nil, fmt.Errorf("invalid unsigned tx data for btc tx, cannot get psbt")
		}
		btcSignTx, err := s.kms.SignTx(seedUUID, params.ChainID, hdPath, btcPsbt)
		if err != nil {
			// this error does not include no row found, need to explicitly check for that
			log.Errorw("failed to sign tx", "err", err)
			return nil, fmt.Errorf("failed to sign tx")
		}
		signedTx = string(btcSignTx)
	} else {
		log.Errorw("invalid chainType")
		return nil, fmt.Errorf("invalid chainID, unknown chainType")
	}
	return &SignTxResponse{
		SignedTx: signedTx,
	}, nil
}
