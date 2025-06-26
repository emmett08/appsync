package main

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/emmett08/dpe-dx-appsync/cmd"
	yamlv3 "gopkg.in/yaml.v3"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	apiServer   *httptest.Server
	apiURL      string
	outputFile  string
	customRegex string
	ref         string
	execErr     error
)

const defaultRegex = `(?P<team>[^_]+)_(?P<owner>[^_]+)_(?P<repo>[^_]+)$`

func aRunningAPIServer(body *godog.DocString) error {
	apiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body.Content)
	}))
	apiURL = apiServer.URL
	return nil
}

func anOutputFilePath(path string) error { outputFile = path; return nil }
func iUseRegex(expr string) error        { customRegex = expr; return nil }
func iUseRef(r string) error             { ref = r; return nil }

func runFetch(owner, repo string) error {
	regex := defaultRegex
	if customRegex != "" {
		regex = customRegex
	}
	args := []string{
		"fetch-repos",
		"--owner", owner, "--repo", repo,
		"--token", "dummy-token",
		"--output", outputFile,
		"--api-url", apiURL,
		"--regex", regex,
	}
	if ref != "" {
		args = append(args, "--ref", ref)
	}
	cmd.RootCmd.SetArgs(args)
	execErr = cmd.RootCmd.Execute()
	return nil
}

func iRunDefault(owner, repo string) error       { return runFetch(owner, repo) }
func iRunWithRegex(owner, repo, rx string) error { customRegex = rx; return runFetch(owner, repo) }
func iRunWithRef(owner, repo, r string) error    { ref = r; return runFetch(owner, repo) }

func theExitStatusShouldBe(code int) error {
	if code == 0 && execErr != nil {
		return fmt.Errorf("expected success, got %v", execErr)
	}
	if code != 0 && execErr == nil {
		return fmt.Errorf("expected exit %d, got 0", code)
	}
	return nil
}

func theFileShouldContain(path string, expect *godog.DocString) error {
	// sanitise path and forbid upward traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid read path: %q", path)
	}
	//nolint:gosec // test helper reading a local sample file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}
	var exp, got interface{}
	if err := yamlv3.Unmarshal([]byte(expect.Content), &exp); err != nil {
		return fmt.Errorf("parse expected YAML: %w", err)
	}
	if err := yamlv3.Unmarshal(data, &got); err != nil {
		return fmt.Errorf("parse actual YAML: %w", err)
	}
	if !reflect.DeepEqual(exp, got) {
		return fmt.Errorf("YAML content mismatch.\n\nwant:\n%s\n\ngot:\n%s", expect.Content, string(data))
	}
	return nil
}

func RegisterFetchReposSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a running API server that returns the following JSON:$`, aRunningAPIServer)
	ctx.Step(`^an output file path "([^"]*)"$`, anOutputFilePath)
	ctx.Step(`^I use regex `+"`"+`(.+)`+"`"+`$`, iUseRegex)
	ctx.Step(`^I use ref "([^"]*)"$`, iUseRef)

	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)"$`, iRunDefault)
	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)" regex "([^"]*)"$`, iRunWithRegex)
	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)" ref "([^"]*)"$`, iRunWithRef)

	ctx.Step(`^the exit status should be (\d+)$`, theExitStatusShouldBe)
	ctx.Step(`^the file "([^"]*)" should contain:$`, theFileShouldContain)
}
