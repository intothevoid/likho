package page

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func CreatePageCmd(cfg *config.Config) *cobra.Command {
	var image, description string

	cmd := &cobra.Command{
		Use:   "page [page-title]",
		Short: "Create a new page",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			err := createPage(cfg, title, image, description)
			if err != nil {
				utils.GetLogger().Error("error creating page: %v\n", zap.Error(err))
				return
			}
		},
	}

	cmd.Flags().StringVarP(&image, "image", "i", "", "Featured image for the page")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Short description of the page")

	return cmd
}

func createPage(cfg *config.Config, title, image, description string) error {
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
		return err
	}

	utils.GetLogger().Sugar().Infof("page created successfully: %s\n", filePath)
	return nil
}
