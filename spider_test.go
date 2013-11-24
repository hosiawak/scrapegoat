package scrapegoat

import (
	"fmt"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type testPost struct {
	site, name, url string
	price           int
	ctxValue        string
}

func newPost() Item {
	return &testPost{site: "amazon.com"}
}

type postCtx struct {
	value string
}

func (p *testPost) Parse(resp *Response, ctx Context) (Item, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	p.name = doc.Find("a").Text()
	p.url = "amazon.com/something"
	p.price, _ = strconv.Atoi(doc.Find("span.price").Text())
	if ctx != nil {
		p.ctxValue = ctx.(*postCtx).value
	}
	return p, nil
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

func TestResponse(t *testing.T) {
	// test server
	body := "<a>Hello</a><span class=price>2</span>"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	// start a spider
	results := make(chan *Response)
	spider := NewSpider("spider.com", results)
	spider.SetLogger(ioutil.Discard)
	spider.NewItemFunc = newPost
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)

	// wait for results
	recv := <-results

	// Response
	if recv.Status != "200 OK" {
		t.Errorf("Expected Status to be '200 OK', got %v", recv.Status)
	}
	if recv.StatusCode != 200 {
		t.Errorf("Expected StatusCode to be 200, got %v", recv.StatusCode)
	}
	if recv.URL.String() != ts.URL {
		t.Errorf("Expected StatusCode to be %s, got %v", recv.URL)
	}

	// Parsed item
	if item, ok := recv.Item.(*testPost); ok {
		if item.site != "amazon.com" {
			t.Errorf("Expected site to be 'amazon.com', got %v", item.site)
		}
		if item.name != "Hello" {
			t.Errorf("Expected name to be 'Hello', got %v", item.name)
		}
		if item.url != "amazon.com/something" {
			t.Errorf("Expected url to be 'amazon.com/something', got %v", item.url)
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

func TestEnqueueURLContext(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	}))
	defer ts.Close()

	// start a spider
	results := make(chan *Response)
	spider := NewSpider("spider.com", results)
	spider.SetLogger(ioutil.Discard)
	spider.NewItemFunc = newPost

	spider.Start()

	// enqueue a URL with Context
	ctx := &postCtx{"value"}

	spider.EnqueueURLContext(ts.URL, ctx)

	// wait for results
	recv := <-results

	// Parsed item
	if item, ok := recv.Item.(*testPost); ok {
		if item.ctxValue != "value" {
			t.Errorf("Expected name to be 'value', got %v", item.ctxValue)
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
	spider.NewItemFunc = newPost

	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)

	// wait for results
	recv := <-c

	if recv.Item.(*testPost).name != "ĄŁŚŹŻ" {
		t.Errorf("Expected to receive 'ĄŁ', got %v", recv)
	}

}
