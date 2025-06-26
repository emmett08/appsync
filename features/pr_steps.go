package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/emmett08/dpe-dx-appsync/internal/app"
)

type FakeRepoGateway struct {
	Default  string
	Branches map[string]string
	Files    map[string]map[string][]byte
	PRs      []PullRequest
}

type PullRequest struct {
	Title string
	Body  string
	Base  string
	Head  string
}

func NewFakeRepo(defaultBranch string) *FakeRepoGateway {
	return &FakeRepoGateway{
		Default:  defaultBranch,
		Branches: map[string]string{defaultBranch: ""},
		Files:    make(map[string]map[string][]byte),
	}
}

func (f *FakeRepoGateway) DefaultBranch(_ context.Context) (string, error) {
	return f.Default, nil
}

func (f *FakeRepoGateway) CreateBranch(_ context.Context, from, to string) error {
	if _, ok := f.Branches[from]; !ok {
		return fmt.Errorf("base branch %s does not exist", from)
	}
	f.Branches[to] = from
	return nil
}

func (f *FakeRepoGateway) WriteFile(_ context.Context, filePath string, content []byte, branch string) error {
	if _, ok := f.Branches[branch]; !ok {
		return fmt.Errorf("branch %s does not exist", branch)
	}
	if f.Files[branch] == nil {
		f.Files[branch] = make(map[string][]byte)
	}
	f.Files[branch][filePath] = content
	return nil
}

func (f *FakeRepoGateway) PullRequest(_ context.Context, title, body, base, head string) (int, error) {
	f.PRs = append(f.PRs, PullRequest{Title: title, Body: body, Base: base, Head: head})
	return len(f.PRs), nil
}

var (
	fakeRepo      *FakeRepoGateway
	selectedStrat app.PRStrategy
	renderedFiles map[string][]byte
	createdBranch string
)

func aRepositoryWithDefaultBranch(def string) error {
	fakeRepo = NewFakeRepo(def)
	return nil
}

func noBranchNamedMatching(pattern string) error {
	re := regexp.MustCompile(pattern)
	for br := range fakeRepo.Branches {
		if re.MatchString(br) {
			return fmt.Errorf("branch %s should not exist", br)
		}
	}
	return nil
}

func theFollowingRenderedFiles(table *godog.Table) error {
	renderedFiles = make(map[string][]byte)
	for _, row := range table.Rows {
		p := row.Cells[0].Value
		renderedFiles[p] = []byte(row.Cells[1].Value)
	}
	return nil
}

func iSelectDirectCommitStrategy() error {
	selectedStrat = app.DirectCommitStrategy{}
	return nil
}

func iApplyTheStrategyToTheRenderedFiles() error {
	return selectedStrat.Apply(context.Background(), fakeRepo, renderedFiles)
}

func theFilesAreWrittenToBranch(branch string) error {
	for p, c := range renderedFiles {
		got, ok := fakeRepo.Files[branch][p]
		if !ok {
			return fmt.Errorf("file %s not written to branch %s", p, branch)
		}
		if !bytes.Equal(got, c) {
			return fmt.Errorf("content mismatch for %s on %s", p, branch)
		}
	}
	return nil
}

func noPullRequestIsCreated() error {
	if len(fakeRepo.PRs) > 0 {
		return fmt.Errorf("expected no PRs, got %d", len(fakeRepo.PRs))
	}
	return nil
}

func iSelectFeatureBranchPRStrategy() error {
	selectedStrat = app.FeatureBranchPRStrategy{}
	return nil
}

func aNewBranchMatchingIsCreatedFrom(pattern, base string) error {
	normalized := strings.ReplaceAll(pattern, `\\`, `\`)
	re, err := regexp.Compile(normalized)
	if err != nil {
		return fmt.Errorf("invalid regex %q: %w", normalized, err)
	}
	for br, b := range fakeRepo.Branches {
		if re.MatchString(br) && b == base {
			createdBranch = br
			return nil
		}
	}
	return fmt.Errorf("no branch matching %q from %q", normalized, base)
}

func theFilesAreWrittenToThatNewBranch() error {
	return theFilesAreWrittenToBranch(createdBranch)
}

func aPullRequestTitledIsOpenedFromTheNewBranchInto(title, base string) error {
	if len(fakeRepo.PRs) != 1 {
		return fmt.Errorf("expected 1 PR, got %d", len(fakeRepo.PRs))
	}
	pr := fakeRepo.PRs[0]
	if pr.Title != title || pr.Base != base || pr.Head != createdBranch {
		return fmt.Errorf("PR mismatch: %+v", pr)
	}
	return nil
}

func RegisterPRSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a repository with default branch "([^"]*)"$`, aRepositoryWithDefaultBranch)
	ctx.Step(`^no branch named matching "([^"]*)"$`, noBranchNamedMatching)
	ctx.Step(`^the following rendered files:$`, theFollowingRenderedFiles)
	ctx.Step(`^I select the direct commit strategy$`, iSelectDirectCommitStrategy)
	ctx.Step(`^I apply the strategy to the rendered files$`, iApplyTheStrategyToTheRenderedFiles)
	ctx.Step(`^the files are written to branch "([^"]*)"$`, theFilesAreWrittenToBranch)
	ctx.Step(`^no pull request is created$`, noPullRequestIsCreated)
	ctx.Step(`^I select the feature-branch PR strategy$`, iSelectFeatureBranchPRStrategy)
	ctx.Step(`^a new branch matching "([^"]*)" is created from "([^"]*)"$`, aNewBranchMatchingIsCreatedFrom)
	ctx.Step(`^the files are written to that new branch$`, theFilesAreWrittenToThatNewBranch)
	ctx.Step(`^a pull request titled "([^"]*)" is opened from the new branch into "([^"]*)"$`, aPullRequestTitledIsOpenedFromTheNewBranchInto)
}
