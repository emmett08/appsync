package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/emmett08/appsync/internal/app"
	"github.com/emmett08/appsync/internal/config"
	"github.com/emmett08/appsync/internal/domain"
)

var (
	tmpRoot     string
	scanner     app.CatalogScanner
	descriptors []domain.ApplicationDescriptor
)

func setupCatalog(paths []string) error {
	var err error
	tmpRoot, err = os.MkdirTemp("", "catalog")
	if err != nil {
		return err
	}
	for _, p := range paths {
		full := filepath.Join(tmpRoot, filepath.FromSlash(p))
		if err := os.MkdirAll(full, fs.ModePerm); err != nil {
			return err
		}
	}
	scanner = app.CatalogScanner{Root: tmpRoot, Filter: config.Filter{}}
	return nil
}

func aCatalogDirectoryWithTheFollowingStructure(table *godog.Table) error {
	var paths []string
	for _, row := range table.Rows {
		paths = append(paths, row.Cells[0].Value)
	}
	return setupCatalog(paths)
}

func iSetFilterTeamTo(team string) error {
	scanner.Filter.Team = team
	return nil
}

func iScanTheCatalogRoot() error {
	var err error
	descriptors, err = scanner.Scan(context.Background())
	return err
}

func iShouldDiscoverDescriptors(count int) error {
	if len(descriptors) != count {
		return fmt.Errorf("expected %d descriptors, got %d", count, len(descriptors))
	}
	return nil
}

func RegisterScanSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a catalog directory with the following structure:$`, aCatalogDirectoryWithTheFollowingStructure)
	ctx.Step(`^I set filter team to "([^"]*)"$`, iSetFilterTeamTo)
	ctx.Step(`^I scan the catalog root$`, iScanTheCatalogRoot)
	ctx.Step(`^I should discover (\d+) descriptors$`, iShouldDiscoverDescriptors)
}
