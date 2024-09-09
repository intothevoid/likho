package post

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/spf13/cobra"
)

func CreatePost(cfg *config.Config) *cobra.Command {
	var tags string
	var featuredImage string

	cmd := &cobra.Command{
		Use:   "create [post-title]",
		Short: "Create a new post",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
			date := time.Now().Format("2006-01-02")

			// Create posts directory if it doesn't exist
			// Create folder structure at root
			postsDir := filepath.Join("../..", cfg.Content.SourceDir, cfg.Content.PostsDir, date)
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
---

Your content here.
`, title, date, strings.Join(tagList, ", "), featuredImage)

			err = os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error writing file: %v\n", err)
				return
			}

			fmt.Printf("Created new post: %s\n", filename)
		},
	}

	cmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags")
	cmd.Flags().StringVarP(&featuredImage, "image", "i", "", "URL of the featured image")

	return cmd
}
