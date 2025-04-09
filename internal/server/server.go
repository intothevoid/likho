package server

import (
	"fmt"
	"net/http"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/generator"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

// Serve starts the HTTP server and serves the generated static files
func Serve(cfg *config.Config) error {
	logger := utils.GetLogger()

	// Generate the static site
	if err := generator.Generate(cfg); err != nil {
		return fmt.Errorf("failed to generate site: %w", err)
	}

	// Set up the file server
	fs := http.FileServer(http.Dir(cfg.Content.OutputDir))
	http.Handle("/", fs)

	// Start the server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server started", zap.String("address", addr))
	return http.ListenAndServe(addr, nil)
}
