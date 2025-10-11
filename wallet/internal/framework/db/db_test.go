package db_test

import (
	"time"

	"wallet/internal/framework/db"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wallet", func() {

	Describe("Configuration", func() {
		DescribeTable("Force you to setup connection resources", func(envOption env.Options) {
			_, err := db.NewConfig(envOption)
			Expect(err).ToNot(BeNil())
		},
			Entry("fail on no env is set", env.Options{
				Environment: map[string]string{},
			}),
			Entry("fail on only env DB_URL is set", env.Options{
				Environment: map[string]string{
					"DB_URL": "postgres://username:password@host:5432/db?sslmode=disable",
				},
			}),
		)

		It("Pass when you know need to setup connection resource", func() {
			envOption := env.Options{
				Environment: map[string]string{
					"DB_URL":            "postgres://username:password@host:5432/db?sslmode=disable",
					"DB_MAX_CONNECTION": "5",
					"DB_MIN_CONNECTION": "1",
					"DB_MAX_IDLE_TIME":  "15m",
				},
			}
			c, err := db.NewConfig(envOption)
			Expect(err).To(BeNil())
			if Expect(c).ToNot(BeNil()) {
				Expect(c.Url).To(Equal("postgres://username:password@host:5432/db?sslmode=disable"))
				Expect(c.MaxConnection).To(Equal(int32(5)))
				Expect(c.MinConnection).To(Equal(int32(1)))
				Expect(c.MaxConnectionIdleTime).To(Equal(15 * time.Minute))
				Expect(c.LogLevel).To(Equal(pgx.LogLevel(pgx.LogLevelError)))
			}
		})

		DescribeTable("Configure log level", func(level string, expected pgx.LogLevel) {
			envOption := env.Options{
				Environment: map[string]string{
					"DB_LOG_LEVEL": level,
					// doesn't matter what values
					"DB_URL":            "1",
					"DB_MAX_CONNECTION": "2",
					"DB_MIN_CONNECTION": "3",
					"DB_MAX_IDLE_TIME":  "4m",
				},
			}
			c, err := db.NewConfig(envOption)
			Expect(err).To(BeNil())
			if Expect(c).ToNot(BeNil()) {
				Expect(c.LogLevel).To(Equal(expected))
			}
		},
			Entry("to info", "info", pgx.LogLevel(pgx.LogLevelInfo)),
			Entry("parse error default to error", "xyz", pgx.LogLevel(pgx.LogLevelError)),
		)
	})

})
