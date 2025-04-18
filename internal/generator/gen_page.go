package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
)

func generatePageHTML(cfg *config.Config, tmpl *template.Template, page parser.Page, pages []parser.Page) error {
	// Convert relative image paths to absolute paths in markdown
	content := page.Content
	content = strings.ReplaceAll(content, "![", "![/images/")
	content = strings.ReplaceAll(content, "](images/", "](/images/")
	content = strings.ReplaceAll(content, "](../images/", "](/images/")
	content = strings.ReplaceAll(content, "](./images/", "](/images/")

	// Convert relative links to files in other directory to absolute paths
	content = strings.ReplaceAll(content, "](other/", "](/other/")
	content = strings.ReplaceAll(content, "](../other/", "](/other/")
	content = strings.ReplaceAll(content, "](./other/", "](/other/")

	data := struct {
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Pages       []parser.Page
	}{
		Content:     template.HTML(content),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   page.Title,
		Pages:       pages,
	}

	// Create pages directory if it doesn't exist
	pagesDir := filepath.Join(cfg.Content.OutputDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create pages directory: %w", err)
	}

	outputPath := filepath.Join(pagesDir, page.Slug+".html")
	return executeTemplate(tmpl, "pages.html", outputPath, data)
}
