package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

// ThemeManager handles theme operations
type ThemeManager struct {
	logger     *zap.Logger
	config     *ThemeConfig
	themePath  string
	outputPath string
}

// NewThemeManager creates a new theme manager
func NewThemeManager(themeName, outputPath string) (*ThemeManager, error) {
	logger := utils.GetLogger()
	themePath := GetThemePath(themeName)

	if !IsValidTheme(themeName) {
		return nil, fmt.Errorf("invalid theme: %s", themeName)
	}

	config, err := LoadThemeConfig(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load theme config: %v", err)
	}

	return &ThemeManager{
		logger:     logger,
		config:     config,
		themePath:  themePath,
		outputPath: outputPath,
	}, nil
}

// CopyAssets copies theme assets to the output directory
func (tm *ThemeManager) CopyAssets() error {
	// Create output directories
	assetDirs := []string{"css", "js", "images"}
	for _, dir := range assetDirs {
		path := filepath.Join(tm.outputPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	// Copy CSS files
	for _, css := range tm.config.Assets.CSS {
		srcPath := filepath.Join(tm.themePath, "static", css)
		dstPath := filepath.Join(tm.outputPath, css)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			tm.logger.Warn("CSS file not found, skipping", zap.String("path", srcPath))
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy CSS file %s: %v", css, err)
		}
	}

	// Copy JS files
	for _, js := range tm.config.Assets.JS {
		srcPath := filepath.Join(tm.themePath, "static", js)
		dstPath := filepath.Join(tm.outputPath, js)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			tm.logger.Warn("JS file not found, skipping", zap.String("path", srcPath))
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy JS file %s: %v", js, err)
		}
	}

	// Copy image files
	for _, img := range tm.config.Assets.Images {
		srcPath := filepath.Join(tm.themePath, "static", img)
		dstPath := filepath.Join(tm.outputPath, img)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			tm.logger.Warn("Image file not found, skipping", zap.String("path", srcPath))
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy image file %s: %v", img, err)
		}
	}

	return nil
}

// GetTemplatePath returns the path to the theme's templates
func (tm *ThemeManager) GetTemplatePath() string {
	return filepath.Join(tm.themePath, "templates")
}

// GetFeatures returns the theme's features
func (tm *ThemeManager) GetFeatures() ThemeFeatures {
	return tm.config.Features
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}
