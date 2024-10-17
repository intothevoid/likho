package generator

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func generateIndexHTML(cfg *config.Config, tmpl *template.Template, posts []post.Post, pages []parser.Page) error {
	data := struct {
		Posts       []post.Post
		Pages       []parser.Page
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		TotalPosts  int
	}{
		Posts:       posts[:min(len(posts), cfg.Content.PostsPerPage)],
		Pages:       pages,
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "Latest",
		TotalPosts:  len(posts),
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "index.html")
	return executeTemplate(tmpl, "index.html", outputPath, data)
}
