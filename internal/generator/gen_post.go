package generator

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func generatePostHTML(cfg *config.Config, tmpl *template.Template, p post.Post, pages []parser.Page) error {
	// Convert Markdown content to HTML with syntax highlighting classes
	extensions := mdparser.CommonExtensions | mdparser.Attributes
	markdownParser := mdparser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	html := markdown.ToHTML([]byte(p.Content), markdownParser, renderer)

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

	// Replace spaces with hyphens and convert to lowercase
	fileName = strings.ToLower(strings.ReplaceAll(fileName, " ", "-"))

	// Replace all non-alphanumeric characters with an empty string
	// Allow hyphens, underscores, and periods
	fileName = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' || r == '.' {
			return r
		}
		return -1
	}, fileName)

	outputPath := filepath.Join(cfg.Content.OutputDir, fileName)

	return executeTemplate(tmpl, "post.html", outputPath, data)
}
