package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
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

	about, err := parser.ParseAboutPage(filepath.Join(cfg.Content.SourceDir, cfg.Content.PagesDir, "about.md"))
	if err != nil {
		return err
	}

	projects, err := parser.ParseProjects(filepath.Join(cfg.Content.SourceDir, cfg.Content.PagesDir, "projects.md"))
	if err != nil {
		return err
	}

	if err := generateHTML(cfg, posts, about, projects); err != nil {
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

func generateHTML(cfg *config.Config, posts []post.Post, about string, projects string) error {
	// Create a FuncMap with custom functions
	funcMap := template.FuncMap{
		"urlize": urlize,
	}

	// Parse all templates with the custom functions
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(cfg.Content.TemplatesDir, "*.html"))
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	// Generate index page
	if err := generateIndexHTML(cfg, tmpl, posts); err != nil {
		return err
	}

	// Generate post pages
	for _, p := range posts {
		if err := generatePostHTML(cfg, tmpl, p); err != nil {
			return err
		}
	}

	// Generate about page
	if err := generateAboutHTML(cfg, tmpl, about); err != nil {
		return err
	}

	// Generate projects page
	if err := generateProjectsHTML(cfg, tmpl, projects); err != nil {
		return err
	}

	// Generate all posts page
	if err := generateAllPostsHTML(cfg, tmpl, posts); err != nil {
		return err
	}

	return nil
}

func generateIndexHTML(cfg *config.Config, tmpl *template.Template, posts []post.Post) error {
	data := struct {
		Posts       []post.Post
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Content     template.HTML
	}{
		Posts:       posts,
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "Home",
		Content:     "",
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "index.html")
	return executeTemplate(tmpl, "index.html", outputPath, data)
}

// urlize converts a string to a URL-friendly format
func urlize(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove any character that isn't a letter, number, or hyphen
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return s
}

func generatePostHTML(cfg *config.Config, tmpl *template.Template, p post.Post) error {
	// Convert Markdown content to HTML
	markdownParser := mdparser.New()
	html := markdown.ToHTML([]byte(p.Content), markdownParser, nil)

	data := struct {
		Post        post.Post
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
	}{
		Post:        p,
		Content:     template.HTML(html),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   p.Title,
	}

	// Create the file name using both the title and the slug
	fileName := fmt.Sprintf("%s-%s.html", p.Title, p.Slug)
	fileName = strings.ToLower(strings.ReplaceAll(fileName, " ", "-"))
	outputPath := filepath.Join(cfg.Content.OutputDir, fileName)

	return executeTemplate(tmpl, "post.html", outputPath, data)
}

func generateAboutHTML(cfg *config.Config, tmpl *template.Template, aboutContent string) error {
	data := struct {
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
	}{
		Content:     template.HTML(aboutContent),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "About",
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "about.html")
	return executeTemplate(tmpl, "about.html", outputPath, data)
}

func generateProjectsHTML(cfg *config.Config, tmpl *template.Template, projectsContent string) error {
	data := struct {
		Posts       []post.Post
		Content     template.HTML
		SiteTitle   string
		CurrentYear int
		PageTitle   string
	}{
		Posts:       nil,
		Content:     template.HTML(projectsContent),
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "Projects",
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "projects.html")
	return executeTemplate(tmpl, "projects.html", outputPath, data)
}

func generateAllPostsHTML(cfg *config.Config, tmpl *template.Template, posts []post.Post) error {
	data := struct {
		Posts       []post.Post
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Content     template.HTML
	}{
		Posts:       posts,
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "All Posts",
		Content:     "", // Leave empty as we're not using it directly
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "posts.html")
	return executeTemplate(tmpl, "posts.html", outputPath, data)
}

func executeTemplate(tmpl *template.Template, name, outputPath string, data interface{}) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", outputPath, err)
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, name, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", name, err)
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
