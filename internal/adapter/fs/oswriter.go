package fs

import (
	"os"
	"path/filepath"

	"github.com/emmett08/1ai-pr/internal/port"
)

type writer struct{}

func New() port.DirWriter { return writer{} }

func (writer) Write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
