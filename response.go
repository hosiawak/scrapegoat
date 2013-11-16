package scrapegoat

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// Response is returned on the spider's results channel
type Response struct {
	r          io.Reader
	Item       interface{} // Holds the processed Item struct (created with NewItemFunc)
	Status     string      // HTTP status, eg. 200 OK
	StatusCode int         // HTTP status code eg. 200
	URL        *url.URL
	Headers    http.Header
	Body       string
	Elapsed    time.Duration // time elapsed between the HTTP request and a response
}
