package scrapegoat

type Item interface {
	// define Item methods here
	Parse(resp *Response, ctx Context) (Item, error)
}
