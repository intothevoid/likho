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

func copyImages(cfg *config.Config) error {
	sourceDir := filepath.Join(cfg.Content.SourceDir, cfg.Content.ImagesDir)
	destinationDir := filepath.Join(cfg.Content.OutputDir, cfg.Content.ImagesDir)

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		// Source directory doesn't exist, create destination directory and return
		return os.MkdirAll(destinationDir, 0755)
	}

	return copyDir(sourceDir, destinationDir)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			err = os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destPath, err)
			}
		} else {
			// Copy file
			sourceFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening source file: %v", err)
			}
			defer sourceFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("error creating destination file: %v", err)
			}
			defer destFile.Close()

			// Preserve permissions
			info, _ := sourceFile.Stat()
			err = os.Chmod(destPath, info.Mode())
			if err != nil {
				return fmt.Errorf("error changing permissions for destination file: %v", err)
			}

			// Copy contents
			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				return fmt.Errorf("error copying file contents: %v", err)
			}

			// Preserve timestamps
			err = os.Chtimes(destPath, info.ModTime(), info.ModTime())
			if err != nil {
				return fmt.Errorf("error setting timestamps for destination file: %v", err)
			}
		}

		return nil
	})
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

func removeImagesDir(outputDir string) error {
	imagesDir := filepath.Join(outputDir, "images")
	// Check if directory exists before trying to remove it
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(imagesDir)
}
