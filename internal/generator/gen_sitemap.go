package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/post"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

func generateSitemap(cfg *config.Config, posts []post.Post) error {
	logger := utils.GetLogger()
	sitemapPath := filepath.Join(cfg.Content.OutputDir, "sitemap.xml")
	file, err := os.Create(sitemapPath)
	if err != nil {
		return fmt.Errorf("error creating sitemap file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	if err != nil {
		return fmt.Errorf("error writing sitemap header: %v", err)
	}

	_, err = file.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	if err != nil {
		return fmt.Errorf("error writing urlset opening tag: %v", err)
	}

	// Add homepage URL
	err = writeURL(file, cfg.Site.BaseURL, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	// Add post URLs
	for _, post := range posts {
		postURL := fmt.Sprintf("%s/%s-%s.html", cfg.Site.BaseURL, urlize(post.Title), post.Slug)
		err = writeURL(file, postURL, post.Date.Format("2006-01-02"))
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString("</urlset>")
	if err != nil {
		return fmt.Errorf("error writing urlset closing tag: %v", err)
	}

	logger.Info("sitemap generated", zap.String("path", sitemapPath))
	return nil
}
