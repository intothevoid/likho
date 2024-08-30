package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/likho/internal/config"
	"github.com/yourusername/likho/internal/generator"
	"github.com/yourusername/likho/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "likho",
		Short: "Likho is a static site generator",
	}

	rootCmd.AddCommand(
		createCmd(cfg),
		generateCmd(cfg),
		serveCmd(cfg),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "create [post-title]",
		Short: "Create a new post",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Implement post creation logic
		},
	}
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
			server.Serve(cfg)
		},
	}
}
