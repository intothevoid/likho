package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Site struct {
		Title       string `mapstructure:"title"`
		Description string `mapstructure:"description"`
		BaseURL     string `mapstructure:"base_url"`
		Language    string `mapstructure:"language"`
	} `mapstructure:"site"`
	Author  string `mapstructure:"author"`
	Content struct {
		SourceDir    string `mapstructure:"source_dir"`
		PostsDir     string `mapstructure:"posts_dir"`
		OutputDir    string `mapstructure:"output_dir"`
		PostsPerPage int    `mapstructure:"posts_per_page"`
	} `mapstructure:"content"`
	Theme struct {
		Name               string `mapstructure:"name"`
		SyntaxHighlighting bool   `mapstructure:"syntax_highlighting"`
	} `mapstructure:"theme"`
	Build struct {
		Draft  bool `mapstructure:"draft"`
		Future bool `mapstructure:"future"`
	} `mapstructure:"build"`
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
	Social struct {
		Twitter  string `mapstructure:"twitter"`
		GitHub   string `mapstructure:"github"`
		LinkedIn string `mapstructure:"linkedin"`
	} `mapstructure:"social"`
	Features struct {
		Comments bool `mapstructure:"comments"`
		Search   bool `mapstructure:"search"`
		RSS      bool `mapstructure:"rss"`
	} `mapstructure:"features"`
	Custom struct {
		GoogleAnalytics string `mapstructure:"google_analytics"`
		DisqusShortname string `mapstructure:"disqus_shortname"`
	} `mapstructure:"custom"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set a default value for PostsDirectory if not specified in config
	if cfg.Content.PostsDir == "" {
		cfg.Content.PostsDir = "posts"
	}

	return &cfg, nil
}
