package scrapegoat

import (
	"net/http"
	"strings"
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
				resp, err := urlReq.send(w.client)

				if err != nil {
					panic(err)
				}

				doc, err := NewDocumentFromReader(strings.NewReader(resp.Body))
				if err != nil {
					return
				}
				item := w.spider.NewItemFunc()
				item.Process(doc, resp, urlReq.ctx)
				resp.Item = item

				w.spider.results <- resp
			case <-w.quit:
				return
			}

		}
	}()
}

func (w *worker) stop() {
	w.quit <- true
}
