scrapegoat
==========

![scrapegoat](https://github.com/hosiawak/scrapegoat/raw/master/logo.png)

Scraper library written in Go.

Usage
==

```go
package main

import (
	"fmt"

	sg "github.com/hosiawak/scrapegoat"
	"github.com/PuerkitoBio/goquery"
)

// Create a struct to hold your parsed item
type Product struct {
	website, name, description string
}

// Create an initialization function for your item
func NewProduct() sg.Item {
	return &Product{website: "amazon.com"}
}

// Define the parsing function Parse for your item
// resp.Body is io.Reader you can use to parse the page
// For example you can use goquery
func (p *Product) Parse(resp *sg.Response, ctx sg.Context) (sg.Item, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	p.name = doc.Find("title").Text()
	p.description = doc.Find(".productDescriptionWrapper").Text()
	return p, nil
}

func main() {

	// Create a channel on which you'll receive *sg.Response:
	results := make(chan *sg.Response)

	// Create a spider and pass the channel
	spider := sg.NewSpider("amazon.com", results)

	// Register the init function at the spider
	spider.NewItemFunc = NewProduct

	// Start the spider
	spider.Start()

	// Enqueue some URLs
	spider.EnqueueURL("http://www.amazon.com/Apple-iPod-classic-Black-Generation/dp/B001F7AHOG")

	// Collect the result
	// This blocks waiting for results so you may want to do it in a goroutine
	res := <-results

	// Need to type assert your struct because on the channel can hold any value (scrapegoat.Item is interface{})
	if product, ok := res.Item.(*Product); ok {
		fmt.Printf("Product Name: %s\nDescription: %s\n", product.name, product.description)
	} else {
		panic("Assertion failed")
	}

	// to stop the spider
	spider.Stop()
}
```
