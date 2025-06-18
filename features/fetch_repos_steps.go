package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
)

var (
	apiServer  *httptest.Server
	apiURL     string
	outputFile string
	cmdErr     error
)

func aRunningAPIServer(body *godog.DocString) error {
	apiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body.Content)
	}))
	apiURL = apiServer.URL
	return nil
}

func anOutputFilePath(p string) error {
	outputFile = p
	return nil
}

func runCLI(extra ...string) {
	base := []string{
		"run", "..", "--", "fetch-repos",
		"--path", "", "--token", "dummy-token",
		"--api-url", apiURL,
		"--output", outputFile,
	}
	args := append(base, extra...)
	cmd := exec.Command("go", args...)
	cmdErr = cmd.Run()
}

func iRunFetchReposDefault(owner, repo string) error {
	runCLI("--owner", owner, "--repo", repo)
	return nil
}

func iRunFetchReposWithRegex(owner, repo, regex string) error {
	runCLI("--owner", owner, "--repo", repo, "--regex", regex)
	return nil
}

func iRunFetchReposWithRef(owner, repo, ref string) error {
	runCLI("--owner", owner, "--repo", repo, "--ref", ref)
	return nil
}

func theExitStatusShouldBe(code int) error {
	if code == 0 && cmdErr == nil {
		return nil
	}
	if ee, ok := cmdErr.(*exec.ExitError); ok && ee.ExitCode() == code {
		return nil
	}
	return fmt.Errorf("expected exit %d, got %v", code, cmdErr)
}

func theFileShouldContain(path string, body *godog.DocString) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), body.Content) {
		return fmt.Errorf("mismatch:\nwant:\n%s\ngot:\n%s", body.Content, string(data))
	}
	return nil
}

func RegisterFetchReposSteps(ctx *godog.ScenarioContext) {
	ctx.After(func(_ context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if apiServer != nil {
			apiServer.Close()
		}
		if outputFile != "" {
			os.Remove(outputFile)
		}
		cmdErr = nil
		return context.Background(), nil
	})

	ctx.Step(`^a running API server that returns the following JSON:$`, aRunningAPIServer)
	ctx.Step(`^an output file path "([^"]*)"$`, anOutputFilePath)
	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)"$`, iRunFetchReposDefault)
	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)" regex "([^"]*)"$`, iRunFetchReposWithRegex)
	ctx.Step(`^I run the fetch-repos command for owner "([^"]*)" repo "([^"]*)" ref "([^"]*)"$`, iRunFetchReposWithRef)
	ctx.Step(`^the exit status should be (\d+)$`, theExitStatusShouldBe)
	ctx.Step(`^the file "([^"]*)" should contain:$`, theFileShouldContain)
}
