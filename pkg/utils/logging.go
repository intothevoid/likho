package utils

import (
	"log"
	"strings"
	"sync"

	"github.com/intothevoid/likho/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes the global logger
func InitLogger(cfg *config.Config) {
	once.Do(func() {
		var level zapcore.Level
		switch strings.ToLower(cfg.Logging.Level) {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		default:
			level = zap.InfoLevel
		}

		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(level)

		var err error
		logger, err = config.Build()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	})
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		log.Fatal("Logger not initialized. Call InitLogger first.")
	}
	return logger
}
