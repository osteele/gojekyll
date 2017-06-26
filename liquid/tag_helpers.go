package liquid

import "github.com/acstech/liquid/core"

type standaloneTag struct {
	name string
}

// AddCode is required by the Liquid tag interface
func (l *standaloneTag) AddCode(code core.Code) {
	panic("AddCode should not have been called on a standalone tag")
}

// AddSibling is required by the Liquid tag interface
func (l *standaloneTag) AddSibling(tag core.Tag) error {
	panic("AddSibling should not have been called on a standalone tag")
}

// LastSibling is required by the Liquid tag interface
func (l *standaloneTag) LastSibling() core.Tag {
	return nil
}

// Name is required by the Liquid tag interface
func (l *standaloneTag) Name() string {
	return l.name
}

// Type is required by the Liquid tag interface
func (l *standaloneTag) Type() core.TagType {
	return core.StandaloneTag
}
