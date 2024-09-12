package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intothevoid/likho/internal/post"
	"gopkg.in/yaml.v2"
	// Add other necessary imports
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
	date, err := time.Parse("2006-01-02", meta.Date)
	if err != nil {
		return post.Post{}, err
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
