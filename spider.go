package scrapegoat

import (
	"io"
	"io/ioutil"
	"log"
)

type Spider struct {
	// The name of the spider
	Name string
	// NewItemFunc is a required func that returns a new item that
	// will be later passed to Process()
	NewItemFunc func() Item
	results     chan *Response // returns responses on this channel
	urlQueue    chan *Request  // URL request queue
	workers     []*worker      // a slice of workers making HTTP requests
	logger      *log.Logger
}

// NewSpider returns a new Spider which will send Responses on the
// supplied channel
func NewSpider(name string, results chan *Response) *Spider {
	spider := &Spider{}
	spider.Name = name
	spider.results = results
	spider.SetLogger(ioutil.Discard)
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
	s.logger.Printf("Setting concurrency to %d\n", c)
	s.workers = make([]*worker, c)
	for i, _ := range s.workers {
		s.workers[i] = s.newWorker()
	}
}

// SetQueueSize sets the size of the buffered URL request channel. It
// is advised to set this before calling EnqueueURL
func (s *Spider) SetQueueSize(b int) {
	s.logger.Printf("Setting queue size to %d\n", b)
	s.urlQueue = make(chan *Request, b)
}

// SetLogger sets the logger for spider
func (s *Spider) SetLogger(w io.Writer) {
	s.logger = log.New(w, "["+s.Name+"] ", log.LstdFlags)
}

// EnqueueURL pushes the given URL onto the URL queue for spidering
func (s *Spider) EnqueueURL(url string) {
	s.logger.Printf("Enqueue URL %s\n", url)
	s.EnqueueURLContext(url, nil)
}

// EnqueueURLContext pushes the given URL and *Context onto the URL queue for
// spidering
func (s *Spider) EnqueueURLContext(url string, ctx Context) {
	req := &Request{url, ctx}
	s.urlQueue <- req
}

// Start starts the spidering process.
func (s *Spider) Start() {
	s.logger.Println("Starting spider")
	for idx, worker := range s.workers {
		s.logger.Printf("Starting worker %d\n", idx)
		worker.start()
	}
}

// Stop stops the spidering process
func (s *Spider) Stop() {
	s.logger.Println("Stopping spider")
	for idx, worker := range s.workers {
		s.logger.Printf("Stopping worker %d\n", idx)
		worker.stop()
	}
}
