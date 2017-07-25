package pages

import (
	"io"
	"os"
)

// A StaticFile is a static file. (Lint made me say this.)
type StaticFile struct {
	file
}

// Static is in the File interface.
func (p *StaticFile) Static() bool { return true }

func (p *StaticFile) Write(w io.Writer) error {
	in, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck, gas
	_, err = io.Copy(w, in)
	return err
}
