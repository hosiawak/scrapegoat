package scrapegoat

import ()

type Spider struct {
	// The name of the spider
	Name string
	// NewItemFunc is a required func that returns a new item that
	// will be later passed to Process()
	NewItemFunc func() Item
	results     chan *Response // returns responses on this channel
	urlQueue    chan *Request  // URL request queue
	workers     []*worker      // a slice of workers making HTTP requests
}

// NewSpider returns a new Spider which will send Responses on the
// supplied channel
func NewSpider(name string, results chan *Response) *Spider {
	spider := &Spider{}
	spider.Name = name
	spider.results = results
	spider.SetConcurrency(2)
	spider.SetQueueSize(100)
	//	debugf("Created spider %s", name)
	return spider
}

// Concurrency returns the number of concurrent HTTP workers active inside
// the spider
func (s *Spider) Concurrency() int {
	return len(s.workers)
}

// SetConcurrency sets the number of concurrent HTTP workers inside
// the spider. It is advised to set the preferred concurrency level
// before calling Start()
func (s *Spider) SetConcurrency(c int) {
	s.workers = make([]*worker, c)
	for i, _ := range s.workers {
		s.workers[i] = s.newWorker()
	}
}

// SetQueueSize sets the size of the buffered URL request channel. It
// is advised to set this before calling EnqueueURL
func (s *Spider) SetQueueSize(b int) {
	s.urlQueue = make(chan *Request, b)
}

// EnqueueURL pushes the given URL onto the URL queue for spidering
func (s *Spider) EnqueueURL(url string) {
	req := &Request{url}
	s.urlQueue <- req
}

// Start starts the spidering process.
func (s *Spider) Start() {
	for _, worker := range s.workers {
		worker.start()
	}
}

// Stop stops the spidering process
func (s *Spider) Stop() {
	for _, worker := range s.workers {
		worker.stop()
	}
}
