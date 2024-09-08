package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/generator"
)

// Serve starts the HTTP server and serves the generated static files
func Serve(cfg *config.Config) error {
	// Generate the static site
	if err := generator.Generate(cfg); err != nil {
		return fmt.Errorf("failed to generate site: %w", err)
	}

	// Set up the file server
	fs := http.FileServer(http.Dir(filepath.Join(".", "public")))
	http.Handle("/", fs)

	// Start the server
	addr := ":8080" // You might want to make this configurable
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
