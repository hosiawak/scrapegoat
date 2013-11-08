package scrapegoat

import (
	"net/http"
	"time"
)

type consumer struct {
	spider         *Spider
	client         *http.Client
	c              chan *http.Response
	requestPending bool
}

func (c *consumer) start() {
	for urlRequest := range c.spider.urlQueue {
		switch urlRequest.action {
		case "fetch":
			req, err := http.NewRequest("GET", urlRequest.url, nil)
			if err != nil {
				panic(err)
			}
			c.requestPending = true
			resp, err := c.client.Do(req)
			if err != nil {
				panic(err)
			}
			c.spider.c <- readResponse(resp)
			resp.Body.Close() // needed to have keep alive
			c.requestPending = false
		case "close":
			ticker := time.Tick(time.Millisecond * 100)
			for _ = range ticker {
				if !c.spider.requestsPending() {
					close(c.spider.c)
				}
			}
		}
	}
}
