package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/intothevoid/likho/internal/post"
)

func TestParsePost(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    post.Post
		wantErr bool
	}{
		{
			name: "Valid post",
			args: args{
				filePath: "test-post.md",
			},
			want: post.Post{
				Title:       "my-test-post",
				Description: "description",
				Date:        time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC),
				Tags:        []string{"test_tag"},
				Content:     "\n\nYour content here.\n",
			},
			wantErr: false,
		},
		{
			name: "Invalid file path",
			args: args{
				filePath: "nonexistent/file.md",
			},
			want:    post.Post{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePost(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Title != tt.want.Title {
				t.Errorf("Title mismatch: got %q, want %q", got.Title, tt.want.Title)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description mismatch: got %q, want %q", got.Description, tt.want.Description)
			}
			if !got.Date.Equal(tt.want.Date) {
				t.Errorf("Date mismatch: got %v, want %v", got.Date, tt.want.Date)
			}
			if !reflect.DeepEqual(got.Tags, tt.want.Tags) {
				t.Errorf("Tags mismatch: got %v, want %v", got.Tags, tt.want.Tags)
			}
			if got.Content != tt.want.Content {
				t.Errorf("Content mismatch:\ngot  %q\nwant %q", got.Content, tt.want.Content)
			}
			// Add this to print the exact content for both got and want
			t.Logf("Got Content: %q", got.Content)
			t.Logf("Want Content: %q", tt.want.Content)
		})
	}
}
