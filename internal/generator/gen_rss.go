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

func generateRSS(cfg *config.Config, posts []post.Post) error {
	logger := utils.GetLogger()
	rssPath := filepath.Join(cfg.Content.OutputDir, "rss.xml")
	file, err := os.Create(rssPath)
	if err != nil {
		return fmt.Errorf("error creating RSS file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	if err != nil {
		return fmt.Errorf("error writing RSS header: %v", err)
	}

	_, err = file.WriteString("<rss version=\"2.0\">\n<channel>\n")
	if err != nil {
		return fmt.Errorf("error writing RSS channel opening tags: %v", err)
	}

	// Write channel information
	_, err = fmt.Fprintf(file, "<title>%s</title>\n", cfg.Site.Title)
	if err != nil {
		return fmt.Errorf("error writing RSS title: %v", err)
	}
	_, err = fmt.Fprintf(file, "<link>%s</link>\n", cfg.Site.BaseURL)
	if err != nil {
		return fmt.Errorf("error writing RSS link: %v", err)
	}
	_, err = fmt.Fprintf(file, "<description>%s</description>\n", cfg.Site.Description)
	if err != nil {
		return fmt.Errorf("error writing RSS description: %v", err)
	}

	// Write items
	for _, post := range posts {
		_, err = file.WriteString("<item>\n")
		if err != nil {
			return fmt.Errorf("error writing item opening tag: %v", err)
		}

		_, err = fmt.Fprintf(file, "<title>%s</title>\n", post.Title)
		if err != nil {
			return fmt.Errorf("error writing item title: %v", err)
		}

		link := fmt.Sprintf("%s/%s-%s.html", cfg.Site.BaseURL, urlize(post.Title), post.Slug)
		_, err = fmt.Fprintf(file, "<link>%s</link>\n", link)
		if err != nil {
			return fmt.Errorf("error writing item link: %v", err)
		}

		_, err = fmt.Fprintf(file, "<pubDate>%s</pubDate>\n", post.Date.Format(time.RFC1123Z))
		if err != nil {
			return fmt.Errorf("error writing item pubDate: %v", err)
		}

		_, err = fmt.Fprintf(file, "<description><![CDATA[%s]]></description>\n", post.Content)
		if err != nil {
			return fmt.Errorf("error writing item description: %v", err)
		}

		_, err = file.WriteString("</item>\n")
		if err != nil {
			return fmt.Errorf("error writing item closing tag: %v", err)
		}
	}

	_, err = file.WriteString("</channel>\n</rss>")
	if err != nil {
		return fmt.Errorf("error writing RSS closing tags: %v", err)
	}

	logger.Info("rss feed generated", zap.String("path", rssPath))
	return nil
}
