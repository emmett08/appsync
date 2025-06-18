package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/emmett08/appsync/internal/config"
	"github.com/emmett08/appsync/internal/domain"
)

type CatalogScanner struct {
	Root   string
	Filter config.Filter
}

func (s CatalogScanner) Scan(ctx context.Context) ([]domain.ApplicationDescriptor, error) {
	var out []domain.ApplicationDescriptor
	err := filepath.WalkDir(s.Root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(s.Root, path)
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) != 2 {
			return nil
		}
		team, app := parts[0], parts[1]
		if !s.Filter.Match(team, app) {
			return filepath.SkipDir
		}
		out = append(out, domain.ApplicationDescriptor{Team: team, App: app, SourcePath: path})
		return filepath.SkipDir
	})
	return out, err
}
