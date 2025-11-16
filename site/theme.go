package site

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (s *Site) findTheme() error {
	if s.cfg.Theme == "" {
		return nil
	}
	exe, err := exec.LookPath("bundle")
	if err != nil {
		return fmt.Errorf("bundle is not in your PATH: %w", err)
	}
	cmd := exec.Command(exe, "show", s.cfg.Theme)
	cmd.Dir = s.AbsDir()
	out, err := cmd.CombinedOutput()
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the %s theme could not be found", s.cfg.Theme)
		}
		return err
	}
	s.themeDir = string(bytes.TrimSpace(out))
	return nil
}

func (s *Site) readThemeAssets() error {
	err := s.readFiles(filepath.Join(s.themeDir, "assets"), s.themeDir)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
