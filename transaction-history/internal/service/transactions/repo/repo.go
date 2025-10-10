package repo

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Db interface {
	pgxscan.Querier
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Repo struct {
	db Db
}

func NewRepo(db Db) *Repo {
	return &Repo{
		db: db,
	}
}
