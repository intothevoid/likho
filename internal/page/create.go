package page

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/spf13/cobra"
)

func CreatePageCmd(cfg *config.Config) *cobra.Command {
	var image, description string

	cmd := &cobra.Command{
		Use:   "page [page-title]",
		Short: "Create a new page",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			createPage(cfg, title, image, description)
		},
	}

	cmd.Flags().StringVarP(&image, "image", "i", "", "Featured image for the page")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Short description of the page")

	return cmd
}

func createPage(cfg *config.Config, title, image, description string) {
	slug := strings.ReplaceAll(strings.ToLower(title), " ", "-")
	fileName := fmt.Sprintf("%s.md", slug)
	filePath := filepath.Join(cfg.Content.SourceDir, cfg.Content.PagesDir, fileName)

	content := fmt.Sprintf(`---
title: "%s"
date: %s
featured_image: "%s"
description: "%s"
---

Write your page content here.
`, title, time.Now().Format("2006-01-02"), image, description)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating page: %v\n", err)
		return
	}

	fmt.Printf("Page created successfully: %s\n", filePath)
}
