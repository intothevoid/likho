package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the configuration for the site
type Config struct {
	Site     SiteConfig     `mapstructure:"site"`
	Content  ContentConfig  `mapstructure:"content"`
	Theme    ThemeConfig    `mapstructure:"theme"`
	Build    BuildConfig    `mapstructure:"build"`
	Server   ServerConfig   `mapstructure:"server"`
	Social   SocialConfig   `mapstructure:"social"`
	Features FeaturesConfig `mapstructure:"features"`
	Custom   CustomConfig   `mapstructure:"custom"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// SiteConfig represents the site configuration
type SiteConfig struct {
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
	BaseURL     string `mapstructure:"base_url"`
	Language    string `mapstructure:"language"`
}

// ContentConfig represents the content configuration
type ContentConfig struct {
	SourceDir    string `mapstructure:"source_dir"`
	PostsDir     string `mapstructure:"posts_dir"`
	OutputDir    string `mapstructure:"output_dir"`
	TemplatesDir string `mapstructure:"templates_dir"`
	PagesDir     string `mapstructure:"pages_dir"`
	PostsPerPage int    `mapstructure:"posts_per_page"`
	ImagesDir    string `mapstructure:"images_dir"`
	OtherDir     string `mapstructure:"other_dir"`
}

// ThemeConfig represents the theme configuration
type ThemeConfig struct {
	Name     string        `mapstructure:"name"`
	Path     string        `mapstructure:"path"`
	Features ThemeFeatures `mapstructure:"features"`
	Custom   ThemeCustom   `mapstructure:"custom"`
}

// ThemeFeatures represents theme-specific features
type ThemeFeatures struct {
	SyntaxHighlighting bool `mapstructure:"syntax_highlighting"`
	DarkMode           bool `mapstructure:"dark_mode"`
}

// ThemeCustom represents custom theme settings
type ThemeCustom struct {
	PrimaryColor string `mapstructure:"primary_color"`
	FontFamily   string `mapstructure:"font_family"`
}

// BuildConfig represents the build configuration
type BuildConfig struct {
	Draft  bool `mapstructure:"draft"`
	Future bool `mapstructure:"future"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// SocialConfig represents the social media configuration
type SocialConfig struct {
	Twitter  string `mapstructure:"twitter"`
	Github   string `mapstructure:"github"`
	Linkedin string `mapstructure:"linkedin"`
}

// FeaturesConfig represents the features configuration
type FeaturesConfig struct {
	Comments bool `mapstructure:"comments"`
	Search   bool `mapstructure:"search"`
	RSS      bool `mapstructure:"rss"`
}

// CustomConfig represents custom configuration
type CustomConfig struct {
	GoogleAnalytics string `mapstructure:"google_analytics"`
	DisqusShortname string `mapstructure:"disqus_shortname"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Set default values for all configuration fields
	// Site defaults
	v.SetDefault("site.title", "My Blog")
	v.SetDefault("site.description", "A blog about technology and programming")
	v.SetDefault("site.base_url", "http://localhost:8080")
	v.SetDefault("site.language", "en")

	// Content defaults
	v.SetDefault("content.source_dir", "content")
	v.SetDefault("content.posts_dir", "posts")
	v.SetDefault("content.output_dir", "public")
	v.SetDefault("content.templates_dir", "templates")
	v.SetDefault("content.pages_dir", "pages")
	v.SetDefault("content.posts_per_page", 10)
	v.SetDefault("content.images_dir", "images")
	v.SetDefault("content.other_dir", "other")

	// Theme defaults
	v.SetDefault("theme.name", "default")
	v.SetDefault("theme.path", "themes/default")
	v.SetDefault("theme.features.syntax_highlighting", true)
	v.SetDefault("theme.features.dark_mode", false)
	v.SetDefault("theme.custom.primary_color", "#2596be")
	v.SetDefault("theme.custom.font_family", "sans-serif")

	// Build defaults
	v.SetDefault("build.draft", false)
	v.SetDefault("build.future", false)

	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "localhost")

	// Social defaults
	v.SetDefault("social.twitter", "")
	v.SetDefault("social.github", "")
	v.SetDefault("social.linkedin", "")

	// Features defaults
	v.SetDefault("features.comments", false)
	v.SetDefault("features.search", true)
	v.SetDefault("features.rss", true)

	// Custom defaults
	v.SetDefault("custom.google_analytics", "")
	v.SetDefault("custom.disqus_shortname", "")

	// Logging defaults
	v.SetDefault("logging.level", "info")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, but that's okay - we'll use defaults
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Debug logging
	fmt.Printf("Loaded config: %+v\n", cfg)
	fmt.Printf("Content section: %+v\n", cfg.Content)

	return &cfg, nil
}
