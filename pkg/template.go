package meetup

import (
	"fmt"
	"os"
	"path"

	"github.com/otiai10/copy"
)

const (
	TemplateDirName = ".templates"
)

// todo: allow set deafult template

func (m *Manager) AddTemplate(paths ...string) error {
	dir := path.Join(m.RootDir, TemplateDirName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not add template: %w", err)
	}

	for _, src := range paths {
		dst := path.Join(dir, path.Base(src))

		if err := copy.Copy(src, dst); err != nil {
			return fmt.Errorf("could not add tmeplate: %w", err)
		}
	}

	return nil
}

func (m *Manager) ListTemplates() ([]string, error) {
	dir := path.Join(m.RootDir, TemplateDirName)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not list templates: %w", err)
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		templates = append(templates, entry.Name())
	}

	return templates, nil
}

func (m *Manager) RemoveTemplate(names ...string) error {
	for _, name := range names {
		if err := os.Remove(path.Join(m.RootDir, TemplateDirName, name)); err != nil {
			return fmt.Errorf("could not remove template: %w", err)
		}
	}

	return nil
}
