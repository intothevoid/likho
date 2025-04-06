package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/internal/parser"
	"github.com/intothevoid/likho/internal/theme"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

func Generate(cfg *config.Config) error {
	logger := utils.GetLogger()

	// Initialize theme manager
	themeManager, err := theme.NewThemeManager(cfg.Theme.Name, cfg.Content.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to initialize theme manager: %v", err)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(cfg.Content.OutputDir, 0755); err != nil {
		logger.Error("error creating output directory", zap.Error(err))
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// remove the images directory
	if err := removeImagesDir(cfg.Content.OutputDir); err != nil {
		return err
	}

	// Remove existing HTML files from the output directory
	if err := removeGeneratedFiles(cfg.Content.OutputDir); err != nil {
		logger.Error("error removing existing HTML files", zap.Error(err))
		return fmt.Errorf("error removing existing HTML files: %v", err)
	}

	posts, err := parser.ParsePosts(filepath.Join(cfg.Content.SourceDir, cfg.Content.PostsDir))
	if err != nil {
		return err
	}

	pages, err := parser.ParsePages(filepath.Join(cfg.Content.SourceDir, cfg.Content.PagesDir))
	if err != nil {
		return err
	}

	// Copy theme assets
	if err := themeManager.CopyAssets(); err != nil {
		return fmt.Errorf("failed to copy theme assets: %v", err)
	}

	// Update templates directory to use theme templates
	cfg.Content.TemplatesDir = themeManager.GetTemplatePath()

	if err := generateHTML(cfg, posts, pages); err != nil {
		return err
	}

	if err := generateTagPages(cfg, posts, pages); err != nil {
		return err
	}

	if err := generateSitemap(cfg, posts); err != nil {
		return err
	}

	if err := generateRSS(cfg, posts); err != nil {
		return err
	}

	if err := copyImages(cfg); err != nil {
		return err
	}

	// Add this summary log at the end of the Generate function
	logger.Info("site generation completed",
		zap.Int("totalPosts", len(posts)),
		zap.Int("totalPages", len(pages)),
		zap.String("outputDir", cfg.Content.OutputDir),
		zap.String("theme", cfg.Theme.Name))

	return nil
}
