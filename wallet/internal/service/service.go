package service

import (
	"context"

	db "wallet/internal/database/generated"
	"wallet/internal/service/utils"
	"wallet/internal/service/validator"

	"github.com/google/uuid"
)

type KMS interface {
	DeriveSeedUUID(accountID uuid.UUID) uuid.UUID
	DerivePublicKeyAndAddress(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents) (string, string, error)
	SignTx(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents, tx any) ([]byte, error)
}

type Repo interface {
	FindWallet(ctx context.Context, accountID string, hdPath utils.HDPathComponents) (*db.Wallet, error)
	SaveWallet(ctx context.Context, wallet db.Wallet) error
}

type Service struct {
	repo Repo
	kms  KMS
	v    *validator.Validator
}

func NewService(kms KMS, repo Repo) *Service {
	return &Service{
		repo: repo,
		kms:  kms,
		v:    validator.NewValidator(),
	}

}
