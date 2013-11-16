package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

// A Document represents the parsed HTML page and allows extraction of
// data from it using CSS selectors
type Document struct {
	r          io.Reader
	goqueryDoc *goquery.Document
}

// NewDocumentFromReader creates a new Document and parses the
// reader's data
func NewDocumentFromReader(r io.Reader) (*Document, error) {
	doc := &Document{r: r}
	gd, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	doc.goqueryDoc = gd
	return doc, nil
}

// CSS accepts a CSS selector and returns a *goquery.Selection
func (d *Document) CSS(selector string) *goquery.Selection {
	return d.goqueryDoc.Find(selector)
}
