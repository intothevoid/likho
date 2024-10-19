package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/intothevoid/likho/internal/post"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v1"
)

func ParsePosts(directory string) ([]post.Post, error) {
	var posts []post.Post

	files, err := filepath.Glob(filepath.Join(directory, "*", "*.md"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		post, err := ParsePost(file)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func ParsePost(filePath string) (post.Post, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.GetLogger().Error("file does not exist", zap.String("filePath", filePath))
		return post.Post{}, fmt.Errorf("file does not exist: %s", filePath)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return post.Post{}, err
	}

	// Split the content into frontmatter and markdown
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return post.Post{}, fmt.Errorf("invalid post format: %s", filePath)
	}

	var meta post.PostMeta
	err = yaml.Unmarshal([]byte(parts[1]), &meta)
	if err != nil {
		return post.Post{}, err
	}

	// Parse the date string into a time.Time object
	date, err := time.Parse("January 2, 2006 15:04", meta.Date)
	if err != nil {
		date, err = time.Parse("2006-01-02", meta.Date)
		if err != nil {
			return post.Post{}, err
		}
	}

	// Create the post
	p := post.Post{
		Title:       meta.Title,
		Description: meta.Description,
		Date:        date,
		Tags:        meta.Tags,
		Content:     parts[2],
		Slug:        filepath.Base(filepath.Dir(filePath)),
	}

	return p, nil
}

type Page struct {
	Title         string
	Date          time.Time
	Description   string
	FeaturedImage string
	Content       string
	Slug          string
}

func ParsePages(directory string) ([]Page, error) {
	var pages []Page

	files, err := filepath.Glob(filepath.Join(directory, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading page file %s: %v", file, err)
		}

		// Split the content into frontmatter and markdown
		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid page format: %s", file)
		}

		var meta struct {
			Title         string    `yaml:"title"`
			Date          time.Time `yaml:"date"`
			FeaturedImage string    `yaml:"featured_image"`
			Description   string    `yaml:"description"`
		}
		err = yaml.Unmarshal([]byte(parts[1]), &meta)
		if err != nil {
			return nil, fmt.Errorf("error parsing frontmatter in %s: %v", file, err)
		}

		// Parse markdown to HTML
		mdParser := parser.New()
		html := markdown.ToHTML([]byte(parts[2]), mdParser, nil)

		slug := strings.TrimSuffix(filepath.Base(file), ".md")

		pages = append(pages, Page{
			Title:         meta.Title,
			Date:          meta.Date,
			FeaturedImage: meta.FeaturedImage,
			Description:   meta.Description,
			Content:       string(html),
			Slug:          slug,
		})
	}

	return pages, nil
}

func ParseAboutPage(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading about page: %v", err)
	}

	// Parse markdown to HTML
	mdParser := parser.New()
	html := markdown.ToHTML(content, mdParser, nil)

	return string(html), nil
}

func ParseProjects(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading projects file: %v", err)
	}

	// Parse markdown to HTML
	mdParser := parser.New()
	html := markdown.ToHTML(content, mdParser, nil)

	return string(html), nil
}

func extractURL(s string) string {
	start := strings.Index(s, "(")
	end := strings.Index(s, ")")
	if start != -1 && end != -1 && start < end {
		return s[start+1 : end]
	}
	return s
}
