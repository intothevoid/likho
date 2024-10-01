package utils

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes the global logger
func InitLogger() {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	})
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}
