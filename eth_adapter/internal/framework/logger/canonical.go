package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// References: https://stripe.com/blog/canonical-log-lines

type Format struct {
	failedMsg      string
	failedLevel    zapcore.Level
	completedMsg   string
	completedLevel zapcore.Level
}

type Option = func(*Format)

func FailedMsg(s string) Option {
	return func(f *Format) {
		f.failedMsg = s
	}
}

func CompletedMsg(s string) Option {
	return func(f *Format) {
		f.completedMsg = s
	}
}

func FailedLevel(l zapcore.Level) Option {
	return func(f *Format) {
		f.failedLevel = l
	}
}

func CompletedLevel(l zapcore.Level) Option {
	return func(f *Format) {
		f.completedLevel = l
	}
}

func (log *Logger) Canonical(start time.Time, err error, options ...Option) {
	defer log.Sync()

	f := &Format{
		failedMsg:      "failed",
		failedLevel:    zapcore.ErrorLevel,
		completedMsg:   "completed",
		completedLevel: zapcore.InfoLevel,
	}

	for _, o := range options {
		o(f)
	}

	log.With("canonical", true, "elapsed", time.Since(start))

	if err != nil {
		log.With("error", err)
		log.Desugar().Log(f.failedLevel, f.failedMsg)
		return
	}

	log.Desugar().Log(f.completedLevel, f.completedMsg)
}
