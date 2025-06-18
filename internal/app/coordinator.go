package app

import (
	"context"
	"path/filepath"
)

type SyncCoordinator struct {
	Scanner    CatalogScanner
	Factory    CRDFactory
	Renderer   ManifestRenderer
	Gateway    RepoGateway
	PRStrategy PRStrategy
	TargetRoot string
}

func (c SyncCoordinator) Sync(ctx context.Context) error {
	descriptors, err := c.Scanner.Scan(ctx)
	if err != nil {
		return err
	}
	for _, d := range descriptors {
		crds := c.Factory.Create(d)
		dest := filepath.Join(c.TargetRoot, d.App)
		files, err := c.Renderer.Render(crds, dest)
		if err != nil {
			return err
		}
		if err := c.PRStrategy.Apply(ctx, c.Gateway, files); err != nil {
			return err
		}
	}
	return nil
}
