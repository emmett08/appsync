package app

import (
	"context"
	"os"
	"path/filepath"

	"github.com/emmett08/appsync/internal/config"
	"github.com/emmett08/appsync/internal/domain"
)

type CatalogScanner struct {
	Root   string
	Filter config.Filter
}

func (s CatalogScanner) Scan(ctx context.Context) ([]domain.ApplicationDescriptor, error) {
	var out []domain.ApplicationDescriptor

	teams, err := os.ReadDir(s.Root)
	if err != nil {
		return nil, err
	}
	for _, teamEntry := range teams {
		if !teamEntry.IsDir() {
			continue
		}
		team := teamEntry.Name()
		if s.Filter.Team != "" && s.Filter.Team != team {
			continue
		}

		teamPath := filepath.Join(s.Root, team)
		apps, err := os.ReadDir(teamPath)
		if err != nil {
			return nil, err
		}
		for _, appEntry := range apps {
			if !appEntry.IsDir() {
				continue
			}
			app := appEntry.Name()
			if s.Filter.App != "" && s.Filter.App != app {
				continue
			}
			out = append(out, domain.ApplicationDescriptor{
				Team:       team,
				App:        app,
				SourcePath: filepath.Join(teamPath, app),
			})
		}
	}

	return out, nil
}
