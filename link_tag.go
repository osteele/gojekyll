package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/acstech/liquid/core"
)

// LinkFactory creates a link tag
func LinkFactory(p *core.Parser, config *core.Configuration) (core.Tag, error) {
	start := p.Position
	p.SkipPastTag()
	end := p.Position - 2
	name := strings.TrimSpace(string(p.Data[start:end]))
	url, ok := site.GetFileURL(name)
	if !ok {
		return nil, fmt.Errorf("link tag: %s not found", name)
	}

	return &Link{url}, nil
}

// Link tag data, for passing information from the factory to Execute
type Link struct {
	url string
}

// AddCode is equired by the Liquid tag interface
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
	if _, err := writer.Write([]byte(l.url)); err != nil {
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
