package post

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/pkg/utils"
)

func TestCreatePost(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := "likho-test"
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatalf("Failed to remove existing temp dir: %v", err)
		}
	}
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		Content: struct {
			SourceDir    string "mapstructure:\"source_dir\""
			PostsDir     string "mapstructure:\"posts_dir\""
			OutputDir    string "mapstructure:\"output_dir\""
			TemplatesDir string "mapstructure:\"templates_dir\""
			PagesDir     string "mapstructure:\"pages_dir\""
			ImagesDir    string "mapstructure:\"images_dir\""
			PostsPerPage int    "mapstructure:\"posts_per_page\""
		}{
			SourceDir:    tempDir,
			PostsDir:     "posts",
			TemplatesDir: "templates",
			PostsPerPage: 10,
		},
	}

	// Init logger
	utils.InitLogger(cfg)

	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedTitle  string
		expectedTags   string
		expectedImage  string
		expectedErrors []string
	}{
		{
			name:          "Basic post creation",
			args:          []string{"Test Post"},
			expectedTitle: "Test Post",
		},
		{
			name:          "Post with tags",
			args:          []string{"Tagged Post"},
			flags:         map[string]string{"tags": "tag1, tag2"},
			expectedTitle: "Tagged Post",
			expectedTags:  "tag1, tag2",
		},
		{
			name:          "Post with featured image",
			args:          []string{"Image Post"},
			flags:         map[string]string{"image": "https://example.com/image.jpg"},
			expectedTitle: "Image Post",
			expectedImage: "https://example.com/image.jpg",
		},
		{
			name:           "No arguments",
			args:           []string{},
			expectedErrors: []string{"accepts 1 arg(s), received 0"},
		},
		{
			name:           "Too many arguments",
			args:           []string{"Too", "Many", "Args"},
			expectedErrors: []string{"accepts 1 arg(s), received 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreatePostCmd(cfg)
			cmd.SetArgs(tt.args)

			for flagName, flagValue := range tt.flags {
				err := cmd.Flags().Set(flagName, flagValue)
				if err != nil {
					t.Fatalf("Failed to set flag %s: %v", flagName, err)
				}
			}

			err := cmd.Execute()

			if len(tt.expectedErrors) > 0 {
				if err == nil {
					t.Errorf("Expected error, got none")
				} else {
					for _, expectedError := range tt.expectedErrors {
						if !strings.Contains(err.Error(), expectedError) {
							t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
						}
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check if the file was created
			expectedSlug := strings.ToLower(strings.ReplaceAll(tt.expectedTitle, " ", "-"))
			currentDate := time.Now().Format("2006-01-02")
			expectedDir := filepath.Join(tempDir, cfg.Content.PostsDir, currentDate)
			expectedFilePath := filepath.Join(expectedDir, expectedSlug+".md")

			// Ensure the directory exists
			err = os.MkdirAll(expectedDir, os.ModePerm)
			if err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			// Wait for a short time to ensure file creation is complete
			time.Sleep(100 * time.Millisecond)

			// Check if the file exists
			if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
				t.Fatalf("Expected file does not exist: %s", expectedFilePath)
			}

			content, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Fatalf("Failed to read created file: %v", err)
			}

			// Check file contents
			if !strings.Contains(string(content), "title: \""+tt.expectedTitle+"\"") {
				t.Errorf("Expected title '%s' not found in file content", tt.expectedTitle)
			}

			if tt.expectedTags != "" && !strings.Contains(string(content), "tags: ["+tt.expectedTags+"]") {
				t.Errorf("Expected tags '%s' not found in file content", tt.expectedTags)
			}

			if tt.expectedImage != "" && !strings.Contains(string(content), "featured_image: \""+tt.expectedImage+"\"") {
				t.Errorf("Expected featured image '%s' not found in file content", tt.expectedImage)
			}
		})
	}
}
