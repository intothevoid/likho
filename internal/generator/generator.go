package generator

import (
	"fmt"
	"html/template"
	"io"
	"log"
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

// Add this new function to remove HTML files
func removeGeneratedFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			if err := os.Remove(path); err != nil {
				return err
			}
			log.Printf("Removed file: %s", path)
		}
		return nil
	})
}

func Generate(cfg *config.Config) error {
	// Remove existing HTML files from the output directory
	if err := removeGeneratedFiles(cfg.Content.OutputDir); err != nil {
		if os.IsNotExist(err) {
			// Create the output directory if it doesn't exist
			if err := os.MkdirAll(cfg.Content.OutputDir, 0755); err != nil {
				log.Printf("error creating output directory: %v", err)
				return fmt.Errorf("error creating output directory: %v", err)
			}
		} else {
			log.Printf("error removing existing HTML files: %v", err)
			return fmt.Errorf("error removing existing HTML files: %v", err)
		}
	}

	posts, err := parser.ParsePosts(filepath.Join(cfg.Content.SourceDir, cfg.Content.PostsDir))
	if err != nil {
		return err
	}

	pages, err := parser.ParsePages(filepath.Join(cfg.Content.SourceDir, cfg.Content.PagesDir))
	if err != nil {
		return err
	}

	if err := generateHTML(cfg, posts, pages); err != nil {
		return err
	}

	if err := generateSitemap(cfg, posts); err != nil {
		return err
	}

	if err := generateRSS(cfg, posts); err != nil {
		return err
	}

	if err := copyCSSFile(cfg); err != nil {
		return err
	}

	return nil
}

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
	log.Printf("Templates parsed: %v", tmpl.DefinedTemplates())

	// Generate index page
	if err := generateIndexHTML(cfg, tmpl, posts, pages); err != nil {
		return err
	}

	// Generate post pages
	tmplPost, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(cfg.Content.TemplatesDir, "base.html"),
		filepath.Join(cfg.Content.TemplatesDir, "post.html"),
		filepath.Join(cfg.Content.TemplatesDir, "header.html"),
		filepath.Join(cfg.Content.TemplatesDir, "pages.html"),
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
			log.Printf("Error generating page %s. Error: %v", page.Title, err)
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
		PageTitle:   "Home",
		TotalPosts:  len(posts),
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "index.html")
	return executeTemplate(tmpl, "index.html", outputPath, data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func generateAllPostsHTML(cfg *config.Config, tmpl *template.Template, posts []post.Post, pages []parser.Page) error {
	data := struct {
		Posts       []post.Post
		SiteTitle   string
		CurrentYear int
		PageTitle   string
		Content     template.HTML
		Pages       []parser.Page
	}{
		Posts:       posts,
		SiteTitle:   cfg.Site.Title,
		CurrentYear: time.Now().Year(),
		PageTitle:   "Posts",
		Content:     "", // Leave empty as we're not using it directly
		Pages:       pages,
	}

	outputPath := filepath.Join(cfg.Content.OutputDir, "posts.html")
	return executeTemplate(tmpl, "posts.html", outputPath, data)
}

func executeTemplate(tmpl *template.Template, name, outputPath string, data interface{}) error {
	log.Printf("Executing template %s for output path %s", name, outputPath)
	log.Printf("Defined templates: %v", tmpl.DefinedTemplates())
	log.Printf("Data being passed to template: %+v", data)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", outputPath, err)
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "base.html", data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", name, err)
	}

	log.Printf("Template %s executed successfully for %s", name, outputPath)
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

func copyCSSFile(cfg *config.Config) error {
	sourcePath := filepath.Join(cfg.Content.SourceDir, "assets", "main.css")
	destPath := filepath.Join(cfg.Content.OutputDir, "main.css")

	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error opening source CSS file: %v", err)
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination CSS file: %v", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error copying CSS file: %v", err)
	}

	log.Printf("CSS file copied to: %s", destPath)
	return nil
}
