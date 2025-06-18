package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/emmett08/appsync/internal/config"
)

type SyncCoordinator struct {
	Scanner    CatalogScanner
	Factory    CRDFactory
	Renderer   ManifestRenderer
	GF         GatewayFactory
	PRStrategy PRStrategy
	TargetRoot string
	Token      string
	Repos      config.RepoConfigs
}

func (c SyncCoordinator) Sync(ctx context.Context) error {
	descriptors, err := c.Scanner.Scan(ctx)
	if err != nil {
		return err
	}
	for _, d := range descriptors {
		owner, repo, ok := c.Repos.ForTeam(d.Team)
		if !ok {
			return fmt.Errorf("no repository configured for team %s", d.Team)
		}
		gw := c.GF.New(c.Token, owner, repo)

		crds := c.Factory.Create(d, owner+"/"+repo)

		dest := filepath.Join(c.TargetRoot, d.App)
		files, err := c.Renderer.Render(crds, dest)
		if err != nil {
			return err
		}

		prefix := filepath.Base(c.Scanner.Root)
		remapped := make(map[string][]byte, len(files))
		for localPath, content := range files {
			rel, err := filepath.Rel(c.TargetRoot, localPath)
			if err != nil {
				return fmt.Errorf("building repo path for %q: %w", localPath, err)
			}
			repoPath := filepath.ToSlash(filepath.Join(prefix, rel))
			remapped[repoPath] = content
		}

		if err := c.PRStrategy.Apply(ctx, gw, remapped); err != nil {
			return err
		}
	}
	return nil
}
