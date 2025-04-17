package generator

import (
	"fmt"
	"html/template"
	"os"
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

	// Convert relative image paths to absolute paths in HTML
	htmlStr := string(html)
	htmlStr = strings.ReplaceAll(htmlStr, "src=\"images/", "src=\"/images/")
	htmlStr = strings.ReplaceAll(htmlStr, "src=\"../images/", "src=\"/images/")
	htmlStr = strings.ReplaceAll(htmlStr, "src=\"./images/", "src=\"/images/")

	// Convert relative links to files in other directory to absolute paths
	htmlStr = strings.ReplaceAll(htmlStr, "href=\"other/", "href=\"/other/")
	htmlStr = strings.ReplaceAll(htmlStr, "href=\"../other/", "href=\"/other/")
	htmlStr = strings.ReplaceAll(htmlStr, "href=\"./other/", "href=\"/other/")

	data := struct {
		Post        post.Post
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Pages       []parser.Page
	}{
		Post:        p,
		Content:     template.HTML(htmlStr),
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

	// Create posts directory if it doesn't exist
	postsDir := filepath.Join(cfg.Content.OutputDir, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return fmt.Errorf("failed to create posts directory: %w", err)
	}

	outputPath := filepath.Join(postsDir, fileName)

	return executeTemplate(tmpl, "post.html", outputPath, data)
}
