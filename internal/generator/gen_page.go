package generator

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
)

func generatePageHTML(cfg *config.Config, tmpl *template.Template, page parser.Page, pages []parser.Page) error {
	data := struct {
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Pages       []parser.Page
	}{
		Content:     template.HTML(page.Content),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   page.Title,
		Pages:       pages,
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, page.Slug+".html")
	return executeTemplate(tmpl, "pages.html", outputPath, data)
}
