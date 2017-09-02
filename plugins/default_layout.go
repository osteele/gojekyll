package plugins

func init() {
	register("jekyll-default-layout", jekyllDefaultLayout{})
}

type jekyllDefaultLayout struct{ plugin }

const (
	pageLayoutName = "page"
	postLayoutName = "post"
	homeLayoutName = "home"
)

var layoutNames = []string{"default", postLayoutName, pageLayoutName, homeLayoutName}

func (p jekyllDefaultLayout) layoutNames(s Site) map[string]string {
	var m = map[string]string{}
	var n string
	for _, k := range layoutNames {
		if s.HasLayout(k) {
			n = k
		}
		m[k] = n
	}
	return m
}

func (p jekyllDefaultLayout) PostInitPage(s Site, pg Page) error {
	fm := pg.FrontMatter()
	if fm["layout"] != nil {
		return nil
	}
	layoutNames := p.layoutNames(s)
	k := pageLayoutName
	switch {
	case pg.IsPost():
		k = postLayoutName
	case pg.URL() == "/":
		k = homeLayoutName
	}
	n := layoutNames[k]
	if n != "" {
		fm["layout"] = n
	}
	return nil
}
