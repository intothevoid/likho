package main

import (
	"log"
	"os"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/generator"
	"github.com/intothevoid/likho/internal/page"
	"github.com/intothevoid/likho/internal/post"
	"github.com/intothevoid/likho/internal/server"
	"github.com/intothevoid/likho/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
		os.Exit(1)
	}

	// Setup logger
	utils.InitLogger(cfg)
	logger := utils.GetLogger()

	rootCmd := &cobra.Command{
		Use:   "likho",
		Short: "Likho is a static site generator",
	}

	rootCmd.AddCommand(createCmd(cfg))
	rootCmd.AddCommand(generateCmd(cfg))
	rootCmd.AddCommand(serveCmd(cfg))

	if err := rootCmd.Execute(); err != nil {
		logger.Error("error executing command", zap.Error(err))
		os.Exit(1)
	}
}

func createCmd(cfg *config.Config) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new post or page",
	}

	createCmd.AddCommand(post.CreatePostCmd(cfg))
	createCmd.AddCommand(page.CreatePageCmd(cfg))

	return createCmd
}

func generateCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate HTML files from markdown",
		Run: func(cmd *cobra.Command, args []string) {
			generator.Generate(cfg)
		},
	}
}

func serveCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the generated blog",
		Run: func(cmd *cobra.Command, args []string) {
			logger := utils.GetLogger()
			logger.Info("Starting server...")
			if err := server.Serve(cfg); err != nil {
				logger.Error("error serving site", zap.Error(err))
				os.Exit(1)
			}
		},
	}
}
