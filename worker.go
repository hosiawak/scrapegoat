package scrapegoat

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type worker struct {
	spider *Spider
	client *http.Client
	quit   chan bool
}

func (s *Spider) newWorker() *worker {
	worker := &worker{}
	worker.spider = s
	worker.client = &http.Client{}
	worker.quit = make(chan bool, 1)
	return worker
}

func (w *worker) start() {
	go func() {

		for {
			select {
			case urlReq := <-w.spider.urlQueue:
				before := time.Now()
				req, err := http.NewRequest("GET", urlReq.url, nil)
				if err != nil {
					panic(err)
				}

				resp, err := w.client.Do(req)
				if err != nil {
					panic(err)
				}

				reader := getReader(resp)

				doc, err := goquery.NewDocumentFromReader(reader)
				if err != nil {
					panic(err)
				}

				w.spider.results <- &Response{Item: w.spider.ItemProcessor(doc),
					Headers: resp.Header,
					Body:    "", // FIXME
					Time:    time.Since(before)}
				resp.Body.Close() // needed to have keep alive
			case <-w.quit:
				return
			}

		}
	}()
}

func (w *worker) stop() {
	w.quit <- true
}
