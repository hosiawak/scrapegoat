package scrapegoat

import "net/http"

type worker struct {
	spider *Spider
	client *http.Client
	quit   chan struct{}
}

func (s *Spider) newWorker() *worker {
	worker := &worker{}
	worker.spider = s
	worker.client = &http.Client{}
	worker.quit = make(chan struct{}, 1)
	return worker
}

func (w *worker) start() {
	go func() {

		for {
			select {
			case urlReq := <-w.spider.urlQueue:
				resp, err := urlReq.send(w.client)

				if err != nil {
					// log the error
					// and continue
					break
				}

				item := w.spider.NewItemFunc()
				item.Parse(resp, urlReq.ctx)
				resp.Item = item

				w.spider.results <- resp
			case <-w.quit:
				return
			}

		}
	}()
}

func (w *worker) stop() {
	w.quit <- struct{}{}
}
