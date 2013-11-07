package scrapegoat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSpider(t *testing.T) {
	spider := NewSpider("spider.com")
	if spider.Concurrency() != 2 {
		t.Errorf("Expected Concurrency to be 2, got %v", spider.Concurrency())
	}
}
func TestSetConcurrency(t *testing.T) {
	spider := NewSpider("spider.com")
	spider.SetConcurrency(3)
	if spider.Concurrency() != 3 {
		t.Errorf("Expected Concurrency to be 3, got %v", spider.Concurrency())
	}
}

func TestEnqueueURL(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	}))
	defer ts.Close()

	// start a spider
	spider := NewSpider("spider.com")
	c := make(chan string)
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL, c)

	// wait for results
	recv := <-c

	if recv != "Hello\n" {
		t.Errorf("Expected to receive 'Hello' on channel, got %v", recv)
	}

}

func TestCharsetConversion(t *testing.T) {
	// test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "charset=iso-8859-2")
		fmt.Fprintln(w, string([]byte{0xA1, 0xA3})) // Ą and Ł in iso-8859-2
	}))
	defer ts.Close()

	// start a spider
	spider := NewSpider("spider.com")
	c := make(chan string)
	spider.Start()

	// enqueue some urls
	spider.EnqueueURL(ts.URL, c)

	// wait for results
	recv := <-c

	if recv != "ĄŁ\n" {
		t.Errorf("Expected to receive 'ĄŁ', got %v", recv)
	}

}
