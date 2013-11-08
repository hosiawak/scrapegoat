package scrapegoat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSpider(t *testing.T) {
	c := make(chan string)
	spider := NewSpider("spider.com", c)
	if spider.Name != "spider.com" {
		t.Errorf("Expected spider.Name to be 'spider.com', got %v", spider.Name)
	}
	if spider.Concurrency() != 2 {
		t.Errorf("Expected Concurrency to be 2, got %v", spider.Concurrency())
	}
}
func TestSetConcurrency(t *testing.T) {
	c := make(chan string)
	spider := NewSpider("spider.com", c)
	spider.SetConcurrency(3)
	if spider.Concurrency() != 3 {
		t.Errorf("Expected Concurrency to be 3, got %v", spider.Concurrency())
	}
}

func TestEnqueueURL(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	}))
	defer ts.Close()

	// start a spider
	c := make(chan string)
	spider := NewSpider("spider.com", c)
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)
	spider.EnqueueURL(ts.URL)

	// stop the spider
	spider.Stop()

	// wait for results
	for recv := range c {
		if recv != "Hello" {
			t.Errorf("Expected to receive 'Hello' on channel, got %v", recv)
		}
	}
}

func TestCharsetConversion(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "charset=iso-8859-2")
		fmt.Fprint(w, string([]byte{0xA1, 0xA3, 0xA6, 0xAC, 0xAF})) // ĄŁŚŹŻ in iso-8859-2
	}))
	defer ts.Close()

	// start a spider
	c := make(chan string)
	spider := NewSpider("spider.com", c)
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL)

	// wait for results
	recv := <-c

	if recv != "ĄŁŚŹŻ" {
		t.Errorf("Expected to receive 'ĄŁ', got %v", recv)
	}

}
