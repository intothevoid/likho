package post

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func CreatePostCmd(cfg *config.Config) *cobra.Command {
	var tags string
	var featuredImage string
	var description string

	cmd := &cobra.Command{
		Use:   "post [post-title]",
		Short: "Create a new post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			title := args[0]
			slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
			date := time.Now()
			createPost(cfg, title, date, slug, tags, featuredImage, description)
			return nil
		},
	}

	cmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags")
	cmd.Flags().StringVarP(&featuredImage, "image", "i", "", "URL of the featured image")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Short description of the post") // Add this line

	return cmd
}

func createPost(cfg *config.Config, title string, date time.Time, slug string, tags string, featuredImage string, description string) {
	// Create posts directory if it doesn't exist
	// Create folder structure at root
	postsDir := filepath.Join(cfg.Content.SourceDir, cfg.Content.PostsDir, date.Format("2006-01-02"))
	if err := os.MkdirAll(postsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create posts directory: %v", err)
	}

	// Generate the filename
	filename := filepath.Join(postsDir, slug+".md")

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Process tags
	tagList := []string{}
	if tags != "" {
		tagList = strings.Split(tags, ",")
		for i, tag := range tagList {
			tagList[i] = strings.TrimSpace(tag)
		}
	}

	content := fmt.Sprintf(`---
title: "%s"
date: %s
draft: true
tags: [%s]
featured_image: "%s"
description: "%s"
---

Your content here.
`, title, date.Format("2006-01-02T15:04:05Z07:00"), strings.Join(tagList, ", "), featuredImage, description)

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		utils.GetLogger().Error("error writing file", zap.Error(err))
		return
	}

	utils.GetLogger().Info("created new post", zap.String("filename", filename))
}
