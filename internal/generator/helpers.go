package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/intothevoid/likho/internal/config"
	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// urlize converts a string to a URL-friendly format
func urlize(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove any character that isn't a letter, number, or hyphen
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return s
}

func copyCSSFile(cfg *config.Config) error {
	sourcePath := filepath.Join("assets", "main.css")
	destPath := filepath.Join(cfg.Content.OutputDir, "main.css")

	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error opening source CSS file: %v", err)
	}
	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination CSS file: %v", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error copying CSS file: %v", err)
	}

	utils.GetLogger().Info("css file copied to", zap.String("path", destPath))
	return nil
}

func writeURL(w io.Writer, loc string, lastmod string) error {
	_, err := fmt.Fprintf(w, "  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n  </url>\n", loc, lastmod)
	if err != nil {
		return fmt.Errorf("error writing URL to sitemap: %v", err)
	}
	return nil
}

func removeGeneratedFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" || filepath.Ext(path) == ".xml" {
			if err := os.Remove(path); err != nil {
				return err
			}
			utils.GetLogger().Info("removed file", zap.String("path", path))
		}
		return nil
	})
}
