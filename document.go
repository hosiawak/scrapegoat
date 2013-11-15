package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type Document struct {
	r          io.Reader
	goqueryDoc *goquery.Document
}

func NewDocumentFromReader(r io.Reader) (*Document, error) {
	doc := &Document{r: r}
	gd, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	doc.goqueryDoc = gd
	return doc, nil
}

func (d *Document) CSS(selector string) *goquery.Selection {
	return d.goqueryDoc.Find(selector)
}
