package scrapegoat

type parser func(css string) string

type ItemLoader struct {
	Fields map[string]parser
}

func (s *Spider) NewItemLoader() *ItemLoader {
	item := &ItemLoader{}
	item.Fields = make(map[string]parser)
	return item
}

func (i *ItemLoader) Field(name string, f parser) {
	i.Fields[name] = f
}
