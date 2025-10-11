package service

import (
	"context"
	"fmt"
	db "wallet/internal/database/generated"
	"wallet/internal/framework/logger"
	"wallet/internal/service/utils"

	"github.com/google/uuid"
)

type GetAddressPublicKeyParams struct {
	AccountId string `json:"accountId" validate:"required"`
	ChainId   string `json:"chainId" validate:"required"`
}

type GetAddressPublicKeyResponse struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
}

func (s *Service) GetAddressPublicKey(ctx context.Context, params *GetAddressPublicKeyParams) (*GetAddressPublicKeyResponse, error) {
	log := logger.GetLogger("GetAddressPublicKey")
	defer func() {
		log.Debugw("completed")
	}()
	if err := s.v.Validate(params); err != nil {
		log.Errorw("invalid params", "err", err)
		return nil, fmt.Errorf("invalid params: %v", err)
	}
	log.With("params", params)
	hdPath, ok := utils.HDPaths[params.ChainId]
	if !ok {
		log.Errorw("invalid chainId, no HDPath")
		return nil, fmt.Errorf("invalid chainId, no hdPath")
	}
	wallet, err := s.repo.FindWallet(ctx, params.AccountId, hdPath)
	if err != nil {
		// this error does not include no row found, need to explicitly check for that
		log.Errorw("failed to find wallet", "err", err)
		return nil, err
	}
	if wallet != nil {
		log.With("foundDB", true, "address", wallet.Address, "publicKey", wallet.PublicKey)
		// already in db, return
		return &GetAddressPublicKeyResponse{
			Address:   wallet.Address,
			PublicKey: wallet.PublicKey,
		}, nil
	}
	log.With("foundDB", false)
	// generate new wallet
	accountId, err := uuid.Parse(params.AccountId)
	if err != nil {
		log.Errorw("invalid accountId", "err", err)
		return nil, fmt.Errorf("invalid accountId")
	}
	seedId := s.kms.DeriveSeedUUID(accountId)
	// get address info
	address, publicKey, err := s.kms.DerivePublicKeyAndAddress(seedId, params.ChainId, hdPath)
	if err != nil {
		log.Errorw("failed to derive public key and address", "err", err)
		return nil, fmt.Errorf("failed to dervice publicKey and address")
	}
	// save in db
	dbWallet := db.Wallet{
		AccountID:    params.AccountId,
		SeedID:       seedId.String(),
		Purpose:      int(hdPath.Purpose),
		CoinType:     int(hdPath.CoinType),
		AccountIndex: int(hdPath.AccountIndex),
		ChangeIndex:  int(hdPath.ChangeIndex),
		AddressIndex: int(hdPath.AccountIndex),
		PublicKey:    publicKey,
		Address:      address,
	}
	// do in parallel to be faster
	go func() {
		err = s.repo.SaveWallet(context.Background(), dbWallet)
		if err != nil {
			log.Errorw("failed to save wallet", "err", err)
		}
	}()
	return &GetAddressPublicKeyResponse{
		Address:   address,
		PublicKey: publicKey,
	}, nil
}
