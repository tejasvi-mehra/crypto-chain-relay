package transactions

import (
	"context"
	"time"
	db "transaction-history/internal/database/generated"
	"transaction-history/internal/service/transactions/repo"
	"transaction-history/internal/service/validator"

	"github.com/caarlos0/env/v6"
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

type Repo interface {
	GetAccountTransactions(ctx context.Context, accountID string, limit, offset int) ([]*db.Txs, error)
	GetAccountBalances(ctx context.Context, accountID string) ([]*repo.AccountBalances, error)
}

type TxService struct {
	repo Repo
	v    *validator.Validator
}

func NewTxService(repo Repo) *TxService {
	return &TxService{
		repo: repo,
		v:    validator.NewValidator(),
	}
}
