package scrapegoat

import (
	"net/http"
	"time"
)

type Response struct {
	Item    interface{}
	Headers http.Header
	Body    string
	Time    time.Duration
}
