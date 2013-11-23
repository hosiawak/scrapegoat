package scrapegoat

type Item interface {
	// define Item methods here
	Process(doc *Document, resp *Response, ctx interface{}) Item
}
