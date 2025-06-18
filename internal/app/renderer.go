package app

import (
	"os"
	"path/filepath"

	"github.com/emmett08/appsync/internal/domain"
)

type ManifestRenderer struct{}

func (ManifestRenderer) Render(crds []domain.CRD, destDir string) (map[string][]byte, error) {
	files := make(map[string][]byte, len(crds))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}
	for _, c := range crds {
		b, err := c.YAML()
		if err != nil {
			return nil, err
		}
		name := filepath.Join(destDir, c.FileName())
		files[name] = b
	}
	return files, nil
}
