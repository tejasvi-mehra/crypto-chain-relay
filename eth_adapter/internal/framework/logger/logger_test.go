package logger_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/caarlos0/env/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"

	"eth-adapter/internal/framework/logger"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Logger", func() {
	It("load from environment", func() {
		config, err := logger.NewLoggerConfig(env.Options{
			Environment: map[string]string{
				"MODE":      "development",
				"LOG_LEVEL": "error",
			},
		})
		if err != nil {
			Fail(fmt.Sprintf("fail to load logger config: %+v", err))
		}

		log := logger.GetLoggerWithConfig("from_environment_1", config.ToZapConfig())
		defer log.Sync()

		log.Info("info logging enabled")
		log.Error("error logging enabled")
		// Output:
		// no output as error only
		// 2022-03-25T14:45:41.509+0800    error   from_environment_1        logger/logger_test.go:49        error logging enabled
		// <<< STACK TRACE HERE >>>

		config, err = logger.NewLoggerConfig(env.Options{
			Environment: map[string]string{
				"MODE":         "development",
				"LOG_LEVEL":    "info",
				"LOG_ENCODING": "json",
			},
		})
		if err != nil {
			Fail(fmt.Sprintf("fail to load logger config: %+v", err))
		}

		log = logger.GetLoggerWithConfig("from_environment_2", config.ToZapConfig())
		defer log.Sync()
		log.Info("info logging enabled")
		// Output:
		// info log is enabled with json now
		// {"level":"info","datetime":"2022-03-25T14:55:46.293+0800","logger":"from_environment_2","caller":"logger/logger_test.go:67","msg":"info logging enabled"}
	})

	It("default is development debug console", func() {
		log := logger.GetLogger("default")
		defer log.Sync()

		log.Info("info logging enabled")
		log.Debug("debug logging enabled")

		// Output:
		// 2022-03-25T14:43:02.964+0800    info    default     logger/logger_test.go:75        info logging enabled
		// 2022-03-25T14:43:02.964+0800    debug   default     logger/logger_test.go:76        debug logging enabled
	})

	It("development debug console", func() {
		config := logger.LoggerConfig{
			Mode:        "development",
			LogLevel:    "debug",
			LogEncoding: "console",
		}

		log := logger.GetLoggerWithConfig("development_debug_console", config.ToZapConfig())
		defer log.Sync()

		log.Info("info logging enabled")
		log.Debug("debug logging enabled")

		// Output:
		// 2022-03-25T14:43:02.964+0800    info    development_debug_console     logger/logger_test.go:92        info logging enabled
		// 2022-03-25T14:43:02.964+0800    debug   development_debug_console     logger/logger_test.go:93        debug logging enabled
	})

	It("production info json", func() {
		config := logger.LoggerConfig{
			Mode:        "production",
			LogLevel:    "info",
			LogEncoding: "json",
		}

		log := logger.GetLoggerWithConfig("production_info_json", config.ToZapConfig())
		defer log.Sync()

		log.Info("info logging enabled")
		log.Debug("debug logging disabled")
		// Output:
		// production is fixed to json
		// {"level":"info","datetime":"2022-03-25T14:43:02.965+0800","logger":"production_info_json","msg":"info logging enabled"}
		// debug is not printed
	})

	It("change log level", func() {
		config := logger.LoggerConfig{
			Mode:        "development",
			LogLevel:    "info",
			LogEncoding: "console",
		}

		log := logger.GetLoggerWithConfig("change_log_level", config.ToZapConfig())
		defer log.Sync()

		log.Info("info logging enabled")
		log.Debug("debug logging enabled")
		// Output:
		// 2022-03-25T14:43:02.964+0800    info    change_log_level     logger/logger_test.go:127        info logging enabled
		// debug is not printed

		log.LogLevel(zapcore.DebugLevel)
		log.Debug("debug logging enabled now")
		// 2022-03-25T14:43:02.964+0800    info    change_log_level     logger/logger_test.go:134        debug logging enabled now
	})

	Describe("canonical", func() {

		sleepTime := time.Duration(1 * time.Second)

		It("completed", func() {
			log := logger.GetLogger("canonical-competed")
			defer log.Canonical(time.Now(), nil)

			time.Sleep(sleepTime)
			// 2022-05-16T14:41:41.502+0800   info    canonical-competed       logger/logger_test.go:136       completed       {"canonical": true, "elapsed": 0.000077416}
		})

		It("completed customisation", func() {
			log := logger.GetLogger("canonical-customisation")
			defer log.Canonical(time.Now(), nil, logger.CompletedMsg("done"), logger.CompletedLevel(zapcore.DebugLevel))

			time.Sleep(sleepTime)
			// 2022-05-16T14:41:41.502+0800   debug   canonical-customisation       logger/logger_test.go:136       done       {"canonical": true, "elapsed": 0.000077416}
		})

		Describe("failed", func() {
			f := func(log *logger.Logger, options ...logger.Option) (err error) {
				defer func(start time.Time) {
					log.Canonical(start, err, options...)
				}(time.Now())

				time.Sleep(sleepTime)
				return errors.New("some error")
			}

			It("standard", func() {
				f(logger.GetLogger("canonical-failed-standard"))
				// 2022-05-16T14:41:41.502+0800   error   canonical-failed       logger/logger_test.go:136       failed       {"canonical": true, "elapsed": 0.000077416, "error": "some error"}
				// stack trace here
			})

			It("in json", func() {
				config := logger.LoggerConfig{
					Mode:        "development",
					LogLevel:    "info",
					LogEncoding: "json",
				}

				f(logger.GetLoggerWithConfig("canonical-failed-in-json", config.ToZapConfig()))
				// {"level":"error","datetime":"2022-05-16T14:41:41.502+0800","logger":"canonical-failed-in-json","caller":"logger/logger_test.go:136","msg":"failed","canonical":true,"elapsed":0.000077416,"error":"some error","stacktrace":"..."}
			})

			It("customisation", func() {
				log := logger.GetLogger("canonical-customisation")
				f(log, logger.FailedMsg("not working"), logger.FailedLevel(zapcore.WarnLevel))
				// 2022-05-16T14:41:41.502+0800   warn   canonical-customisation       logger/logger_test.go:136       not working       {"canonical": true, "elapsed": 0.000077416, "error": "some error"}
				// stack trace here
			})
		})

	})

})
