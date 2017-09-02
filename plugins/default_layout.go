package plugins

func init() {
	register("jekyll-default-layout", jekyllDefaultLayout{})
}

type jekyllDefaultLayout struct{ plugin }

const (
	defaultLayout int = iota
	pageLayout
	postLayout
	homeLayout
)

var layoutNames = []string{"default", "post", "page", "home"}

func (p jekyllDefaultLayout) layoutNames(s Site) []string {
	var ln string
	names := make([]string, len(layoutNames))
	for i, n := range layoutNames {
		// fmt.Println("examine", i, n, names)
		if s.HasLayout(n) {
			ln = n
		}
		names[i] = ln
	}
	return names
}

func (p jekyllDefaultLayout) PostInitPage(s Site, pg Page) error {
	fm := pg.FrontMatter()
	if fm["layout"] != nil {
		return nil
	}
	layoutNames := p.layoutNames(s)
	ln := layoutNames[pageLayout]
	switch {
	case pg.IsPost():
		ln = layoutNames[postLayout]
	case pg.Permalink() == "/":
		ln = layoutNames[homeLayout]
	}
	if ln != "" {
		fm["layout"] = ln
	}
	return nil
}
