package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/gomarkdown/markdown"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func Generate(cfg *config.Config) error {
	posts, err := parser.ParsePosts(filepath.Join(cfg.Content.SourceDir, cfg.Content.PostsDir))
	if err != nil {
		return err
	}

	if err := generateHTML(cfg, posts); err != nil {
		return err
	}

	if err := generateSitemap(cfg, posts); err != nil {
		return err
	}

	if err := generateRSS(cfg, posts); err != nil {
		return err
	}

	return nil
}

func generateHTML(cfg *config.Config, posts []post.Post) error {

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(cfg.Content.OutputDir, 0755); err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "post.html"))
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	for _, p := range posts {
		// Convert Markdown content to HTML
		markdownParser := mdparser.New()
		html := markdown.ToHTML([]byte(p.Content), markdownParser, nil)

		data := struct {
			Post        post.Post
			Content     template.HTML
			SiteTitle   string
			CurrentYear int
		}{
			Post:        p,
			Content:     template.HTML(html),
			SiteTitle:   cfg.Site.Title,
			CurrentYear: time.Now().Year(),
		}

		// Create a new file for the HTML output
		outputPath := filepath.Join(cfg.Content.OutputDir, p.Slug+".html")
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tmpl.Execute(file, data)
		if err != nil {
			return fmt.Errorf("error executing template: %v", err)
		}
	}

	return nil
}

func generateSitemap(cfg *config.Config, posts []post.Post) error {
	// Implement sitemap generation logic
	return nil
}

func generateRSS(cfg *config.Config, posts []post.Post) error {
	// Implement RSS feed generation logic
	return nil
}
