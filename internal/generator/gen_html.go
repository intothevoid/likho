package generator

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

func generateHTML(cfg *config.Config, posts []post.Post, pages []parser.Page) error {
	// Create a FuncMap with custom functions
	funcMap := template.FuncMap{
		"urlize": urlize,
	}

	// Parse all templates with the custom functions
	// tmpl, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(cfg.Content.TemplatesDir, "*.html"))
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "index.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "footer.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}
	utils.GetLogger().Debug("templates parsed", zap.Int("numTemplates", len(tmpl.DefinedTemplates())))

	// Generate index page
	if err := generateIndexHTML(cfg, tmpl, posts, pages); err != nil {
		return err
	}

	// Generate post pages
	tmplPost, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "post.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "footer.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}
	for _, p := range posts {
		if err := generatePostHTML(cfg, tmplPost, p, pages); err != nil {
			return err
		}
	}

	// Generate pages
	tmpPages, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "pages.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "footer.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	// Generate html for all pages
	for _, page := range pages {
		if err := generatePageHTML(cfg, tmpPages, page, pages); err != nil {
			utils.GetLogger().Error("error generating page", zap.String("title", page.Title), zap.Error(err))
			return err
		}
	}

	// Generate all posts page
	tmplPosts, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "posts.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "footer.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}
	if err := generateAllPostsHTML(cfg, tmplPosts, posts, pages); err != nil {
		return err
	}

	return nil
}
