package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/generator"
	"github.com/intothevoid/likho/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")      // Look in current directory as well
	viper.AddConfigPath("../../") // Look two directories up from cmd/likho
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	viper.Set("config", &cfg)

	rootCmd := &cobra.Command{
		Use:   "likho",
		Short: "Likho is a static site generator",
	}

	confPtr, ok := viper.Get("config").(*config.Config)
	if !ok {
		log.Fatal("Failed to retrieve configuration")
	}

	rootCmd.AddCommand(createCmd(confPtr))
	rootCmd.AddCommand(generateCmd(confPtr))
	rootCmd.AddCommand(serveCmd(confPtr))

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
			title := args[0]
			slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
			date := time.Now().Format("2006-01-02")

			postDir := filepath.Join(cfg.SourceDirectory, date)
			err := os.MkdirAll(postDir, 0755)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
				return
			}

			fileName := filepath.Join(postDir, slug+".md")
			content := fmt.Sprintf(`---
title: "%s"
date: %s
draft: true
---

Your content here.
`, title, date)

			err = os.WriteFile(fileName, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error writing file: %v\n", err)
				return
			}

			fmt.Printf("Created new post: %s\n", fileName)
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
