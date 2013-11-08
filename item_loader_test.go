package scrapegoat

import (
	"testing"
)

func TestField(t *testing.T) {
	c := make(chan string)
	spider := NewSpider("test", c)
	post := spider.NewItemLoader()
	post.Field("test", func(a string) string { return "ok" })
	if len(post.Fields) != 1 {
		t.Errorf("Expected item.Fields to contain 1 field, got %v", len(post.Fields))
	}
	if post.Fields["test"]("a") != "ok" {
		t.Errorf("Expected item.Field['test']() to return 'ok', got %v", post.Fields["test"]("a"))
	}
}
