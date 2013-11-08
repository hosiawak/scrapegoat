package scrapegoat

import (
	"net/http"
)

type Spider struct {
	Name      string
	c         chan string // returns responses on this channel
	urlQueue  chan *urlRequest
	consumers []*consumer
}

type urlRequest struct {
	action string
	url    string // fetch this url
}

func (s *Spider) newConsumer() *consumer {
	consumer := &consumer{}
	consumer.spider = s
	consumer.client = &http.Client{}
	consumer.c = make(chan *http.Response)
	return consumer
}

func NewSpider(name string, c chan string) *Spider {
	spider := &Spider{}
	spider.Name = name
	spider.c = c
	spider.SetConcurrency(2)
	spider.SetBufferSize(100)
//	debugf("Created spider %s", name)
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

func (s *Spider) EnqueueURL(url string) {
	req := &urlRequest{"fetch", url}
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

func (s *Spider) Stop() {
	s.urlQueue <- &urlRequest{"close", ""}
}

func (s *Spider) requestsPending() bool {
	for _, consumer := range s.consumers {
		if consumer.requestPending {
			return true
		}
	}
	return false
}
