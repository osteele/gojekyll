package collection

import (
	"github.com/osteele/gojekyll/pages"
)

type pagesByDate struct{ pages []pages.Page }

// Len is part of sort.Interface.
func (p pagesByDate) Len() int {
	return len(p.pages)
}

// Less is part of sort.Interface.
func (p pagesByDate) Less(i, j int) bool {
	a, b := p.pages[i].PostDate(), p.pages[j].PostDate()
	return a.After(b)
}

// Swap is part of sort.Interface.
func (p pagesByDate) Swap(i, j int) {
	pages := p.pages
	pages[i], pages[j] = pages[j], pages[i]
}
