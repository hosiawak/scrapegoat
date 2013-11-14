package scrapegoat

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type post struct {
	name  string
	price int
}

func TestNewSpider(t *testing.T) {
	c := make(chan *Response)
	spider := NewSpider("spider.com", c)
	if spider.Name != "spider.com" {
		t.Errorf("Expected spider.Name to be 'spider.com', got %v", spider.Name)
	}
	if spider.Concurrency() != 2 {
		t.Errorf("Expected Concurrency to be 2, got %v", spider.Concurrency())
	}
}

func TestSetConcurrency(t *testing.T) {
	c := make(chan *Response)
	spider := NewSpider("spider.com", c)
	spider.SetConcurrency(3)
	if spider.Concurrency() != 3 {
		t.Errorf("Expected Concurrency to be 3, got %v", spider.Concurrency())
	}
}

func TestEnqueueURL(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<a>Hello</a><span class=price>2</span>")
	}))
	defer ts.Close()

	// start a spider
	results := make(chan *Response)
	spider := NewSpider("spider.com", results)
	spider.ItemProcessor = func(doc *goquery.Document) Item {
		post := &post{}
		post.name = doc.Find("a").Text()
		post.price, _ = strconv.Atoi(doc.Find("span.price").Text())
		return post
	}
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)

	// wait for results
	recv := <-results

	if item, ok := recv.Item.(*post); ok {
		if item.name != "Hello" {
			t.Errorf("Expected name to be 'Hello', got %v", item.name)
		}
		if item.price != 2 {
			t.Errorf("Expected price to be 2, got %v", item.price)
		}
	} else {
		t.Error("Assertion failed")
	}
	// stop the spider
	spider.Stop()

}

func TestCharsetConversion(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "charset=iso-8859-2")
		fmt.Fprint(w, "<a>"+string([]byte{0xA1, 0xA3, 0xA6, 0xAC, 0xAF})+"</a>") // ĄŁŚŹŻ in iso-8859-2
	}))
	defer ts.Close()

	// start a spider
	c := make(chan *Response)
	spider := NewSpider("spider.com", c)
	spider.ItemProcessor = func(doc *goquery.Document) Item {
		post := &post{}
		post.name = doc.Find("a").Text()
		return post
	}
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)

	// wait for results
	recv := <-c

	if recv.Item.(*post).name != "ĄŁŚŹŻ" {
		t.Errorf("Expected to receive 'ĄŁ', got %v", recv)
	}

}
