package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/emmett08/1ai-pr/internal/adapter/git"
	"github.com/emmett08/1ai-pr/internal/adapter/github"
	"github.com/emmett08/1ai-pr/internal/adapter/logging"
	"github.com/emmett08/1ai-pr/internal/app/rollout"
	"github.com/emmett08/1ai-pr/internal/shared/clock"
)

/* -------------------------------------------------------------------------- */
/*  Configuration                                                             */
/* -------------------------------------------------------------------------- */

type config struct {
	Skeleton        string   `yaml:"skeleton"`
	ApplicationsDir string   `yaml:"applicationsDir"`
	TenantRepos     []string `yaml:"tenantRepos"`
}

/* -------------------------------------------------------------------------- */
/*  CLI                                                                       */
/* -------------------------------------------------------------------------- */

func main() {
	root := &cobra.Command{Use: "1ai-pr"}

	var cfgPath string
	var dryRun bool

	rolloutCmd := &cobra.Command{
		Use:   "rollout",
		Short: "Roll out skeleton into tenant repos",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgPath == "" {
				return fmt.Errorf("--config is required")
			}

			if dryRun {
				return runDryRun(cfgPath)
			}

			ctx := context.Background()
			gitRepo := logging.DecorateRepo(git.New())
			prSvc := logging.DecoratePR(github.New(os.Getenv("GITHUB_TOKEN")))
			clk := clock.System{}

			uc := rollout.NewUseCase(gitRepo, prSvc, clk)
			return uc.Execute(ctx, cfgPath)
		},
	}

	rolloutCmd.Flags().StringVar(&cfgPath, "config", "", "path to rollout config")
	rolloutCmd.Flags().BoolVar(&dryRun, "dry-run", false, "copy files locally without Git/PR")
	root.AddCommand(rolloutCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

/* -------------------------------------------------------------------------- */
/*  Dry-run helper (used by the tests)                                        */
/* -------------------------------------------------------------------------- */

func runDryRun(cfgPath string) error {
	// Work with an absolute path for determinism.
	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return err
	}
	cfgDir := filepath.Dir(absCfgPath)

	// -----------------------------------------------------------------------
	// Load and parse the YAML file.
	// -----------------------------------------------------------------------
	raw, err := os.ReadFile(absCfgPath)
	if err != nil {
		return err
	}

	var cfg config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return err
	}

	// Use ".applications" unless overridden.
	applicationsDir := cfg.ApplicationsDir
	if applicationsDir == "" {
		applicationsDir = ".applications"
	}

	// -----------------------------------------------------------------------
	// Resolve the skeleton directory — first try the path “as-is”
	// (relative to the current working directory); if that fails, fall back
	// to being relative to the config file’s directory.
	// -----------------------------------------------------------------------
	skeletonDir := cfg.Skeleton
	if !filepath.IsAbs(skeletonDir) {
		if _, err := os.Stat(skeletonDir); os.IsNotExist(err) {
			skeletonDir = filepath.Join(cfgDir, skeletonDir)
		}
	}
	if _, err := os.Stat(skeletonDir); err != nil {
		return fmt.Errorf("skeleton directory %q: %w", skeletonDir, err)
	}

	skelEntries, err := os.ReadDir(skeletonDir)
	if err != nil {
		return err
	}

	// -----------------------------------------------------------------------
	// Process each tenant repository.
	// -----------------------------------------------------------------------
	for _, rp := range cfg.TenantRepos {
		repoPath := rp
		if !filepath.IsAbs(repoPath) {
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				repoPath = filepath.Join(cfgDir, repoPath)
			}
		}

		for _, entry := range skelEntries {
			if !entry.IsDir() {
				continue
			}
			appName := entry.Name()
			srcAppDir := filepath.Join(skeletonDir, appName)
			destAppDir := filepath.Join(repoPath, applicationsDir, appName)

			if err := copyDir(srcAppDir, destAppDir); err != nil {
				return err
			}
		}

		// Fake PR URL for the dry-run mode (used by the BDD tests).
		fmt.Printf(
			"https://example.com/%s/pull/1 branch 1ai/\n",
			filepath.Base(repoPath),
		)
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*  Utility                                                                   */
/* -------------------------------------------------------------------------- */

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		srcFile := filepath.Join(src, e.Name())
		dstFile := filepath.Join(dst, e.Name())

		data, err := os.ReadFile(srcFile)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dstFile, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
