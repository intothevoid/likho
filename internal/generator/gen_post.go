package generator

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func generatePostHTML(cfg *config.Config, tmpl *template.Template, p post.Post, pages []parser.Page) error {
	// Convert Markdown content to HTML
	markdownParser := mdparser.New()
	html := markdown.ToHTML([]byte(p.Content), markdownParser, nil)

	data := struct {
		Post        post.Post
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Pages       []parser.Page
	}{
		Post:        p,
		Content:     template.HTML(html),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   p.Title,
		Pages:       pages,
	}

	// Create the file name using both the title and the slug
	fileName := fmt.Sprintf("%s-%s.html", p.Title, p.Slug)
	fileName = strings.ToLower(strings.ReplaceAll(fileName, " ", "-"))
	outputPath := filepath.Join(cfg.Content.OutputDir, fileName)

	return executeTemplate(tmpl, "post.html", outputPath, data)
}
