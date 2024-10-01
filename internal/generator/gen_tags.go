package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func generateTagPages(cfg *config.Config, posts []post.Post, pages []parser.Page) error {
	// Create a FuncMap with custom functions
	funcMap := template.FuncMap{
		"urlize": urlize,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "tags.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "footer.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	tags := make(map[string][]post.Post)
	for _, p := range posts {
		for _, tag := range p.Tags {
			tags[tag] = append(tags[tag], p)
		}
	}

	for tag, tagPosts := range tags {
		data := struct {
			Posts       []post.Post
			Pages       []parser.Page
			SiteTitle   string
			CurrentYear int
			PageTitle   string
			Tag         string
		}{
			Posts:       tagPosts,
			Pages:       pages,
			SiteTitle:   cfg.Site.Title,
			CurrentYear: time.Now().Year(),
			PageTitle:   fmt.Sprintf("Posts tagged with %s", tag),
			Tag:         tag,
		}

		// Use urlize function here to ensure consistency
		outputPath := filepath.Join(cfg.Content.OutputDir, "tags", urlize(tag)+".html")
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}
		if err := executeTemplate(tmpl, "tags.html", outputPath, data); err != nil {
			return err
		}
	}

	return nil
}
