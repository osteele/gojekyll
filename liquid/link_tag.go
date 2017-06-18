package liquid

import (
	"fmt"
	"io"
	"strings"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

func init() {
	liquid.Tags["link"] = LinkFactory
}

// A LinkHandler given an include tag file name returns a URL.
type LinkHandler func(string) (string, bool)

var currentLinkHandler LinkHandler

// SetLinkHandler sets the function that resolves an include tag file name to a URL.
func SetLinkHandler(h LinkHandler) {
	currentLinkHandler = h
}

// LinkFactory creates a link tag
func LinkFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	start := p.Position
	p.SkipPastTag()
	end := p.Position - 2
	filename := strings.TrimSpace(string(p.Data[start:end]))
	return &Link{filename}, nil
}

// Link tag data, for passing information from the factory to Execute
type Link struct {
	filename string
}

// AddCode is required by the Liquid tag interface
func (l *Link) AddCode(code core.Code) {
	panic("AddCode should not have been called on a Link")
}

// AddSibling is required by the Liquid tag interface
func (l *Link) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a Link")
}

// LastSibling is required by the Liquid tag interface
func (l *Link) LastSibling() core.Tag {
	return nil
}

// Execute is required by the Liquid tag interface
func (l *Link) Execute(writer io.Writer, data map[string]interface{}) core.ExecuteState {
	url, ok := currentLinkHandler(l.filename)
	if !ok {
		panic(fmt.Errorf("link tag: %s not found", l.filename))
	}
	if _, err := writer.Write([]byte(url)); err != nil {
		panic(err)
	}
	return core.Normal
}

// Name is required by the Liquid tag interface
func (l *Link) Name() string {
	return "link"
}

// Type is required by the Liquid tag interface
func (l *Link) Type() core.TagType {
	return core.StandaloneTag
}
