package scrapegoat

import (
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Item       interface{}
	Status     string
	StatusCode int
	URL        *url.URL
	Headers    http.Header
	Body       string
	Elapsed    time.Duration
}
