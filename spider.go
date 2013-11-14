package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
)

type Spider struct {
	Name          string
	ItemProcessor func(*goquery.Document) Item
	results       chan *Response // returns responses on this channel
	urlQueue      chan *urlRequest
	workers       []*worker
}

type urlRequest struct {
	url string // fetch this url
}

func NewSpider(name string, results chan *Response) *Spider {
	spider := &Spider{}
	spider.Name = name
	spider.results = results
	spider.SetConcurrency(2)
	spider.SetQueueSize(100)
	//	debugf("Created spider %s", name)
	return spider
}

func (s *Spider) Concurrency() int {
	return len(s.workers)
}

func (s *Spider) SetConcurrency(c int) {
	s.workers = make([]*worker, c)
	for i, _ := range s.workers {
		s.workers[i] = s.newWorker()
	}
}
func (s *Spider) SetQueueSize(b int) {
	s.urlQueue = make(chan *urlRequest, b)
}

func (s *Spider) EnqueueURL(url string) {
	req := &urlRequest{url}
	s.urlQueue <- req
}

func (s *Spider) Start() {
	for _, worker := range s.workers {
		worker.start()
	}
}

func (s *Spider) Stop() {
	for _, worker := range s.workers {
		worker.stop()
	}
}
