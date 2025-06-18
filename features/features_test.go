package main

import (
	"testing"

	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	RegisterFetchReposSteps(ctx)
	RegisterScanSteps(ctx)
	RegisterPRSteps(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "appsync",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("feature tests failed")
	}
}
