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
)

// Create a struct to hold your parsed item
type Product struct {
	website, name, description string
}

// Create an initialization function for your item
func NewProduct() sg.Item {
	return &Product{website: "amazon.com"}
}

// Define the parsing function Process for your item
func (p *Product) Process(doc *sg.Document, resp *sg.Response, ctx interface{}) sg.Item {
	p.name = doc.CSS("title").Text()
	p.description = doc.CSS(".productDescriptionWrapper").Text()
	return p
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
