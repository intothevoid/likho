package generator

import (
	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/post"
)

func Generate(cfg *config.Config) error {
	posts, err := parser.ParsePosts(cfg.Content.PostsDir)
	if err != nil {
		return err
	}

	if err := generateHTML(cfg, posts); err != nil {
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

func generateHTML(cfg *config.Config, posts []post.Post) error {
	// Implement HTML generation logic
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
