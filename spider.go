package scrapegoat

import (
	"net/http"
)

type Spider struct {
	urlQueue  chan *urlRequest
	consumers []*consumer
}

type urlRequest struct {
	url string      // fetch this url
	c   chan string // return response on this channel
}

func (s *Spider) newConsumer() *consumer {
	consumer := &consumer{}
	consumer.spider = s
	consumer.client = &http.Client{}
	consumer.c = make(chan *http.Response)
	return consumer
}

func NewSpider(name string) *Spider {
	spider := &Spider{}
	spider.SetConcurrency(2)
	spider.SetBufferSize(100)
	return spider
}

func (s *Spider) Concurrency() int {
	return len(s.consumers)
}

func (s *Spider) SetConcurrency(c int) {
	s.consumers = make([]*consumer, c)
	for i, _ := range s.consumers {
		s.consumers[i] = s.newConsumer()
	}
}
func (s *Spider) SetBufferSize(b int) {
	s.urlQueue = make(chan *urlRequest, b)
}

func (s *Spider) EnqueueURL(url string, c chan string) {
	req := &urlRequest{url, c}
	s.urlQueue <- req
}

func (s *Spider) Start() {
	go s.startConsumers()
}

func (s *Spider) startConsumers() {
	for _, consumer := range s.consumers {
		go consumer.start()
	}
}
