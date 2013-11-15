package scrapegoat

import ()

type Spider struct {
	Name          string
	NewItem       func() Item
	results       chan *Response // returns responses on this channel
	urlQueue      chan *Request
	workers       []*worker
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
	s.urlQueue = make(chan *Request, b)
}

func (s *Spider) EnqueueURL(url string) {
	req := &Request{url}
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
