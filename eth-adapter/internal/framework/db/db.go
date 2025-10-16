package db

import (
	"context"
	"reflect"
	"strings"
	"time"

	"eth-adapter/internal/framework/logger"

	"github.com/caarlos0/env/v6"
	zapadapter "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type Config struct {
	Url                   string            `env:"DB_URL,required"`
	MaxConnection         int32             `env:"DB_MAX_CONNECTION,required"`
	MinConnection         int32             `env:"DB_MIN_CONNECTION,required"`
	MaxConnectionIdleTime time.Duration     `env:"DB_MAX_IDLE_TIME,required"`
	LogLevel              tracelog.LogLevel `env:"DB_LOG_LEVEL" envDefault:"ERROR"`
}

type Db struct {
	*pgxpool.Pool
}

func NewConfig(options ...env.Options) (*Config, error) {
	c := &Config{}
	var logLevel tracelog.LogLevel
	if err := env.ParseWithFuncs(c,
		map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(logLevel): func(v string) (interface{}, error) {
				l, err := tracelog.LogLevelFromString(strings.ToLower(v))
				if err != nil {
					l = tracelog.LogLevelError
				}
				return l, nil
			},
		}, options...); err != nil {
		return nil, err
	}
	return c, nil
}

func NewDb(config *Config) (*Db, error) {
	pgConfig, err := pgxpool.ParseConfig(config.Url)
	if err != nil {
		return nil, err
	}

	pgConfig.MaxConns = config.MaxConnection
	pgConfig.MinConns = config.MinConnection
	pgConfig.MaxConnIdleTime = config.MaxConnectionIdleTime
	pgConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		LogLevel: config.LogLevel,
		Logger:   zapadapter.NewLogger(logger.GetLogger("db").Desugar()),
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}

	return &Db{
		Pool: pool,
	}, nil
}
