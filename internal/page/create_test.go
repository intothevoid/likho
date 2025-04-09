package page

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreatePageCmd(t *testing.T) {
	cfg := &config.Config{
		Content: config.ContentConfig{
			SourceDir: "testdata",
			PagesDir:  "pages",
		},
	}

	cmd := CreatePageCmd(cfg)

	if cmd.Use != "page [page-title]" {
		t.Errorf("Expected Use to be 'page [page-title]', got %s", cmd.Use)
	}

	if cmd.Short != "Create a new page" {
		t.Errorf("Expected Short to be 'Create a new page', got %s", cmd.Short)
	}

	// Check the Args expectation
	if cmd.Args == nil {
		t.Errorf("Expected Args to be set, but it was nil")
	} else if cmd.Args(cmd, []string{"test"}) != nil {
		t.Errorf("Expected Args to accept exactly one argument")
	} else if cmd.Args(cmd, []string{}) == nil {
		t.Errorf("Expected Args to reject zero arguments")
	} else if cmd.Args(cmd, []string{"test1", "test2"}) == nil {
		t.Errorf("Expected Args to reject more than one argument")
	}

	// Ensure Run is set
	assert.NotNil(t, cmd.Run)

	flags := cmd.Flags()
	imageFlag := flags.Lookup("image")
	if imageFlag == nil || imageFlag.Shorthand != "i" {
		t.Errorf("Expected image flag with shorthand 'i'")
	}

	descFlag := flags.Lookup("description")
	if descFlag == nil || descFlag.Shorthand != "d" {
		t.Errorf("Expected description flag with shorthand 'd'")
	}
}

func TestCreatePage(t *testing.T) {
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

	cfg := &config.Config{
		Content: config.ContentConfig{
			SourceDir: tempDir,
			PagesDir:  "pages",
		},
	}

	// Init logger
	utils.InitLogger(cfg)

	title := "test-page"
	image := "test-image.jpg"
	description := "This is a test page"

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(filepath.Join(tempDir, "pages", title+".md")), os.ModePerm)
	assert.NoError(t, err)

	err = createPage(cfg, title, image, description)
	assert.NoError(t, err)

	// Check if the file was created
	expectedFilePath := filepath.Join(tempDir, "pages", title+".md")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created, but it doesn't exist", expectedFilePath)
	}

	// Read the content of the created file
	content, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	// Check if the content is correct
	expectedContent := `---
title: "test-page"
date: `
	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("Expected content to contain %s, but it doesn't", expectedContent)
	}

	if !strings.Contains(string(content), "featured_image: \"test-image.jpg\"") {
		t.Errorf("Expected content to contain featured_image: \"test-image.jpg\", but it doesn't")
	}

	if !strings.Contains(string(content), "description: \"This is a test page\"") {
		t.Errorf("Expected content to contain description: \"This is a test page\", but it doesn't")
	}

	if !strings.Contains(string(content), "Write your page content here.") {
		t.Errorf("Expected content to contain 'Write your page content here.', but it doesn't")
	}
}
