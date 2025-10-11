package service

import (
	"context"
	"time"

	db "wallet/internal/database/generated"
	"wallet/internal/service/utils"
	"wallet/internal/service/validator"

	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
)

type ServiceConfig struct {
	Host         string        `env:"HOST" envDefault:"127.0.0.1"`
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
}

func NewServiceConfig() (*ServiceConfig, error) {
	c := &ServiceConfig{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

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
