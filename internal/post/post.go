package post

import (
	"time"
)

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Tags        []string
	Content     string
	Slug        string
}

type PostMeta struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
}
