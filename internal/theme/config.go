package theme

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ThemeConfig represents the configuration for a theme
type ThemeConfig struct {
	Name        string        `yaml:"name"`
	Version     string        `yaml:"version"`
	Description string        `yaml:"description"`
	Author      string        `yaml:"author"`
	License     string        `yaml:"license"`
	Assets      ThemeAssets   `yaml:"assets"`
	Features    ThemeFeatures `yaml:"features"`
}

// ThemeAssets represents the assets included in a theme
type ThemeAssets struct {
	CSS    []string `yaml:"css"`
	JS     []string `yaml:"js"`
	Images []string `yaml:"images"`
}

// ThemeFeatures represents the features supported by a theme
type ThemeFeatures struct {
	SyntaxHighlighting bool `yaml:"syntax_highlighting"`
	Responsive         bool `yaml:"responsive"`
	DarkMode           bool `yaml:"dark_mode"`
}

// LoadThemeConfig loads the theme configuration from the specified path
func LoadThemeConfig(themePath string) (*ThemeConfig, error) {
	v := viper.New()
	v.SetConfigName("theme")
	v.SetConfigType("yaml")
	v.AddConfigPath(themePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config ThemeConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Convert relative paths to absolute paths
	for i, css := range config.Assets.CSS {
		// Remove "static/" prefix if it exists
		css = strings.TrimPrefix(css, "static/")
		config.Assets.CSS[i] = css
	}
	for i, js := range config.Assets.JS {
		// Remove "static/" prefix if it exists
		js = strings.TrimPrefix(js, "static/")
		config.Assets.JS[i] = js
	}
	for i, img := range config.Assets.Images {
		// Remove "static/" prefix if it exists
		img = strings.TrimPrefix(img, "static/")
		config.Assets.Images[i] = img
	}

	return &config, nil
}

// GetThemePath returns the absolute path to the theme directory
func GetThemePath(themeName string) string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Join("themes", themeName)
	}
	return filepath.Join(cwd, "themes", themeName)
}

// GetThemeAssetPath returns the absolute path to a theme asset
func GetThemeAssetPath(themeName, assetPath string) string {
	themePath := GetThemePath(themeName)
	return filepath.Join(themePath, assetPath)
}

// IsValidTheme checks if a theme is valid by verifying its configuration file
func IsValidTheme(themeName string) bool {
	themePath := GetThemePath(themeName)
	configPath := filepath.Join(themePath, "theme.yaml")

	// Check if theme directory and config file exist
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	return true
}
