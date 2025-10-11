
package logger

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger

	// Use this to change log level dynamically
	level zap.AtomicLevel
}

func (log *Logger) LogLevel(newLevel zapcore.Level) {
	log.level.SetLevel(newLevel)
}

func (log *Logger) With(args ...interface{}) {
	log.SugaredLogger = log.SugaredLogger.With(args...)
}

type LoggerConfig struct {
	Mode        string `env:"MODE" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	LogEncoding string `env:"LOG_ENCODING" envDefault:"console"`
}

func NewLoggerConfig(options ...env.Options) (*LoggerConfig, error) {
	c := &LoggerConfig{}
	if err := env.Parse(c, options...); err != nil {
		return nil, err
	}
	return c, nil
}

func GetDefaultZapConfig() zap.Config {
	if lc, err := NewLoggerConfig(); err != nil {
		panic(fmt.Errorf("fail to load logger config: %+v", err))
	} else {
		return lc.ToZapConfig()
	}
}

func (c *LoggerConfig) ToZapConfig() zap.Config {
	var logConfig zap.Config
	if c.Mode == "production" {
		logConfig = zap.NewProductionConfig()
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.EncoderConfig.TimeKey = "datetime"
		logConfig.DisableCaller = true
		logConfig.Encoding = "json"
	} else {
		logConfig = zap.NewDevelopmentConfig()
		logConfig.EncoderConfig = zap.NewProductionEncoderConfig()
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.EncoderConfig.TimeKey = "datetime"
		logConfig.Encoding = c.LogEncoding
	}

	atom := zap.NewAtomicLevel()
	if c.LogLevel != "" {
		newLevel := zap.NewAtomicLevel()
		err := newLevel.UnmarshalText([]byte(c.LogLevel))
		if err != nil {
			fmt.Printf("ignored to invalid log level '%s'\n", c.LogLevel)
		} else {
			atom.SetLevel(newLevel.Level())
		}
	}
	logConfig.Level = atom

	return logConfig
}

func GetLogger(name string, options ...zap.Option) *Logger {
	return GetLoggerWithConfig(name, GetDefaultZapConfig(), options...)
}

func GetLoggerWithConfig(name string, logConfig zap.Config, options ...zap.Option) *Logger {
	log, err := logConfig.Build(options...)
	if err != nil {
		panic("cannot get logger is dead!")
	}

	return &Logger{log.Sugar().Named(name), logConfig.Level}
}
