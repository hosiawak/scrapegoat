package scrapegoat

import (
	"bytes"
	"net/http"
	"time"
)

// Request represents a spidering request. You don't instantiate it
// directly but by calling EnqueueURL or EnqueueRequest (not
// implemented yet).
type Request struct {
	url string // fetch this url
}

// sends the HTTP request, receives a response and encodes it using
// charsetReader as UTF-8
func (r *Request) send(c *http.Client) (*Response, error) {
	before := time.Now()
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	defer resp.Body.Close() // needed to have keep alive

	if err != nil {
		return nil, err
	}
	reader := charsetReader(resp)
	var body bytes.Buffer
	_, _ = body.ReadFrom(reader)

	if err != nil {
		return nil, err
	}

	response := &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL,
		Headers:    resp.Header,
		Body:       body.String(),
		Elapsed:    time.Since(before)}
	return response, nil
}
