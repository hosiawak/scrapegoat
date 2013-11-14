scrapegoat
==========

![scrapegoat](https://github.com/hosiawak/scrapegoat/raw/master/logo.png)

Scraper library written in Go.

Usage
==

    import (
	       "fmt"
           sg "github.com/hosiawak/scrapegoat"
    )
     
    func main() {

         // Create a struct to hold your parsed item
		 type Product { website, name string, price float64 }
		 
         // Create a channel on which you'll receive parsed results:
         results := make(chan *sg.Response)

         // Create a spider and pass the channel
		 spider := sg.NewSpider("amazon.com", results)

         // Create a parsing function ItemProcessor on the spider
		 spider.ItemProcessor = func(doc *goquery.Document) sg.Item {
			 product := &Product{}
			 product.website = "amazon.com"
			 product.name = doc.Find("#btAsInTitle").Text()
			 product.price = doc.Find("b.priceLarge.kitsunePrice").Text()
			 return post
         }

         // Start the spider
		 spider.Start()

         // Enqueue some URLs
		 spider.EnqueueURL("http://www.amazon.com/gp/product/B00CTUKFNQ")
		 spider.EnqueueURL("http://www.amazon.com/Apple-iPod-classic-Black-Generation/dp/B001F7AHOG")
		 // Collect the result
		 // This blocks so you may want to do it in a goroutine
		 res := <-results

         // Need to type assert your struct because on the channel can hold any value (scrapegoat.Item is interface{})
		 if product, ok := res.Item.(*Product); ok {
               fmt.Printf("Product name is %s, price is %d, product.name, product.price)
         } else {
			 panic("Assertion failed")
		 }

		 // to stop the spider
		 spider.Stop()
