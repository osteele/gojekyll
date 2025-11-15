package pages

import (
	"io"
	"os"
)

// A StaticFile is a static file. (Lint made me say this.)
type StaticFile struct {
	file
}

// IsStatic is in the File interface.
func (d *StaticFile) IsStatic() bool { return true }

func (d *StaticFile) Write(w io.Writer) error {
	in, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	defer in.Close() // nolint: errcheck
	_, err = io.Copy(w, in)
	return err
}
