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
				w.spider.logger.Printf("Sending request to %s\n", urlReq.url)

				resp, err := urlReq.send(w.client)
				w.spider.logger.Printf("Received response from %s in %v\n", urlReq.url, resp.Elapsed)

				if err != nil {
					// log the error
					w.spider.logger.Printf("Received error %s", err)
					// and continue
					break
				}

				item := w.spider.NewItemFunc()
				w.spider.logger.Println("Parsing item")
				item.Parse(resp, urlReq.ctx)
				w.spider.logger.Printf("Parsed item %+v\n", item)
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
