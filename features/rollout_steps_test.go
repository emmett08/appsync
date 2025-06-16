package features

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

var (
	lastStdout bytes.Buffer
	lastStderr bytes.Buffer
	lastErr    error
)

func theSkeletonDirectory(path string, table *godog.Table) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		rel := row.Cells[0].Value
		content := row.Cells[1].Value
		full := filepath.Join(path, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func aFakeTenantRepoAt(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}

	init := exec.Command("git", "init", "--initial-branch=main")
	init.Dir = path
	if err := init.Run(); err != nil {
		return fmt.Errorf("git init: %w", err)
	}

	cfgName := exec.Command("git", "config", "user.name", "Test User")
	cfgName.Dir = path
	if err := cfgName.Run(); err != nil {
		return fmt.Errorf("git config user.name: %w", err)
	}
	cfgEmail := exec.Command("git", "config", "user.email", "test@example.com")
	cfgEmail.Dir = path
	if err := cfgEmail.Run(); err != nil {
		return fmt.Errorf("git config user.email: %w", err)
	}

	commit := exec.Command("git", "commit", "--allow-empty", "-m", "initial")
	commit.Dir = path
	if err := commit.Run(); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	return nil
}

func iExecuteCLI(cmdline string) error {
	parts := strings.Fields(cmdline)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Env = append(os.Environ(), "GITHUB_TOKEN=fake")
	lastStdout.Reset()
	lastStderr.Reset()
	cmd.Stdout = &lastStdout
	cmd.Stderr = &lastStderr
	lastErr = cmd.Run()
	return nil
}

func directoryShouldContain(dir string, table *godog.Table) error {
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		expect := row.Cells[0].Value
		full := filepath.Join(dir, expect)
		if _, err := os.Stat(full); err != nil {
			return fmt.Errorf("missing %q in %q: %w", expect, dir, err)
		}
	}
	return nil
}

func aPullRequestShouldHaveBeenOpenedFor(repo, branchPrefix string) error {
	out := lastStdout.String()
	if !strings.Contains(out, "https://example.com/"+repo+"/pull/") {
		return fmt.Errorf("no PR URL for %q", repo)
	}
	if !strings.Contains(out, "branch "+branchPrefix) {
		return fmt.Errorf("expected branch %q", branchPrefix)
	}
	return nil
}

func aConfigFile(path string, table *godog.Table) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("cannot create dir for %q: %w", path, err)
	}

	if len(table.Rows) < 2 {
		return fmt.Errorf("config table must have at least 2 rows, got %d", len(table.Rows))
	}

	skeleton := table.Rows[0].Cells[1].Value
	tenantRepo := table.Rows[1].Cells[1].Value

	cfg := map[string]interface{}{
		"skeleton":    skeleton,
		"tenantRepos": []string{tenantRepo},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a skeleton directory "([^"]*)" containing:$`, theSkeletonDirectory)
	ctx.Step(`^a fake tenant repo at "([^"]*)"(?: with a main branch)?$`, aFakeTenantRepoAt)
	ctx.Step(`^a config file "([^"]*)" containing:$`, aConfigFile)
	ctx.Step(`^I execute `+"`([^`]*)`"+`$`, iExecuteCLI)
	ctx.Step(`^the directory "([^"]*)" should contain:$`, directoryShouldContain)
	ctx.Step(`^a pull request should have been opened for "([^"]*)" on branch "([^"]*)"$`, aPullRequestShouldHaveBeenOpenedFor)
}

func TestFeatures(t *testing.T) {
	status := godog.TestSuite{
		Name:                "rollout",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"."},
			Randomize: time.Now().UTC().UnixNano(),
		},
	}.Run()

	if status != 0 {
		t.Fatalf("godog suite failed with status %d", status)
	}
}
