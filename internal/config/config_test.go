package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file for testing
	configContent := `
content:
  source_dir: "content"
  posts_dir: ""
  templates_dir: ""
`
	configFile, err := os.Create("config.yaml")
	assert.NoError(t, err)
	defer os.Remove(configFile.Name())

	_, err = configFile.Write([]byte(configContent))
	assert.NoError(t, err)
	configFile.Close()

	viper.SetConfigFile(configFile.Name())
	viper.AddConfigPath(".")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "content", cfg.Content.SourceDir)
	assert.Equal(t, "posts", cfg.Content.PostsDir)
	assert.Equal(t, "templates", cfg.Content.TemplatesDir)
}
