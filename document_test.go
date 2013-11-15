package scrapegoat

import (
	"strings"
	"testing"
)

func TestNewDocumentFromReader(t *testing.T) {
	r := strings.NewReader("a")
	doc, _ := NewDocumentFromReader(r)

	if doc.r != r {
		t.Errorf("Expected Document.r to be %v, got %v", r, doc.r)
	}
}

func TestCSSExtract(t *testing.T) {
	r := strings.NewReader("<html><body><a>Amazon</a></body></html>")
	doc, _ := NewDocumentFromReader(r)

	e := doc.CSS("a").Text()

	if e != "Amazon" {
		t.Errorf("Expected to get 'Amazon', got %v", e)
	}
}
