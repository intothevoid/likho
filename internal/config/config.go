package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Title          string
	URL            string
	Language       string
	Description    string
	DateFormat     string
	Author         string
	FrontPagePosts int
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

	return &cfg, nil
}
