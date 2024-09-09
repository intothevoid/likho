package post

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/intothevoid/likho/internal/config"
)

func TestCreatePost(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := "likho-test"
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(filepath.Join("../..", tempDir))

	// Create a test config
	cfg := &config.Config{
		Content: struct {
			SourceDir    string "mapstructure:\"source_dir\""
			PostsDir     string "mapstructure:\"posts_dir\""
			OutputDir    string "mapstructure:\"output_dir\""
			PostsPerPage int    "mapstructure:\"posts_per_page\""
		}{
			SourceDir: tempDir,
			PostsDir:  "posts",
		},
	}

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
			cmd := CreatePost(cfg)
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
			expectedDir := filepath.Join("../..", tempDir, cfg.Content.PostsDir, currentDate)
			expectedFilePath := filepath.Join(expectedDir, expectedSlug+".md")

			// Wait for a short time to ensure file creation is complete
			time.Sleep(100 * time.Millisecond)

			// Print the expected file path
			t.Logf("Expected file path: %s", expectedFilePath)

			// Check if the directory exists
			if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
				t.Fatalf("Expected directory does not exist: %s", expectedDir)
			}

			// List files in the directory
			files, err := os.ReadDir(expectedDir)
			if err != nil {
				t.Fatalf("Failed to read directory: %v", err)
			}

			// Print out all files in the directory
			t.Logf("Files in directory %s:", expectedDir)
			for _, file := range files {
				t.Logf("- %s", file.Name())
			}

			// Check if the expected file exists
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
