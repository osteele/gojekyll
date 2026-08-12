package renderers

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

type templateRoot struct {
	path string
	root *os.Root
}

type rootedTemplateStore struct {
	roots []templateRoot
}

func newRootedTemplateStore(paths ...string) (*rootedTemplateStore, error) {
	store := &rootedTemplateStore{}
	for _, path := range paths {
		if path == "" {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		root, err := os.OpenRoot(abs)
		if err != nil {
			store.close()
			return nil, err
		}
		store.roots = append(store.roots, templateRoot{path: abs, root: root})
	}
	runtime.AddCleanup(store, func(roots []templateRoot) {
		for _, root := range roots {
			_ = root.root.Close()
		}
	}, store.roots)
	return store, nil
}

func (s *rootedTemplateStore) ReadTemplate(filename string) ([]byte, error) {
	abs := filename
	if !filepath.IsAbs(abs) {
		var err error
		abs, err = filepath.Abs(abs)
		if err != nil {
			return nil, err
		}
	}
	for _, root := range s.roots {
		rel, err := filepath.Rel(root.path, filepath.Clean(abs))
		if err != nil || !filepath.IsLocal(rel) {
			continue
		}
		contents, err := root.root.ReadFile(rel)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return contents, err
		}
	}
	return nil, &os.PathError{Op: "open", Path: filename, Err: fs.ErrNotExist}
}

func (s *rootedTemplateStore) close() {
	for _, root := range s.roots {
		_ = root.root.Close()
	}
}
