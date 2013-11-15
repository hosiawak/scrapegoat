package scrapegoat

import (
	"bytes"
	"net/http"
	"time"
)

type Request struct {
	url string // fetch this url
}

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
	reader := getReader(resp)
	var body bytes.Buffer
	_, _ = body.ReadFrom(reader)

	response := &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL,
		Headers:    resp.Header,
		Body:       body.String(), // FIXME
		Elapsed:    time.Since(before)}
	return response, nil
}
