package scrapegoat

import (
	"net/http"
)

type consumer struct {
	spider *Spider
	client *http.Client
	c      chan *http.Response
}

func (c *consumer) start() {
	for {
		urlRequest := <-c.spider.urlQueue
		req, err := http.NewRequest("GET", urlRequest.url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := c.client.Do(req)
		if err != nil {
			panic(err)
		}
		urlRequest.c <- readResponse(resp)
		resp.Body.Close() // needed to have keep alive
	}
}
