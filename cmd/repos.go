package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reposCmd)
	reposCmd.AddCommand(reposStatusCmd)
	reposCmd.AddCommand(reposPullCmd)
	reposStatusCmd.Flags().StringArrayVarP(&repoDirs, "dir", "d", nil, "directories to scan (default: cwd)")
	reposStatusCmd.Flags().BoolVar(&reposNoFetch, "no-fetch", false, "skip git fetch (faster, may be stale)")
	reposPullCmd.Flags().StringArrayVarP(&repoDirs, "dir", "d", nil, "directories to scan (default: cwd)")
}

var repoDirs []string
var reposNoFetch bool

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Local repository utilities",
}

var reposPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fast-forward pull all repos (skips diverged)",
	Run:   runReposPull,
}

var reposStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show ahead/behind status for all local git repos",
	Run:   runReposStatus,
}

type repoInfo struct {
	path   string
	name   string
	branch string
	ahead  int
	behind int
	dirty  bool
	err    string
}

func runReposStatus(cmd *cobra.Command, args []string) {
	dirs := repoDirs
	if len(dirs) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot get cwd:", err)
			os.Exit(1)
		}
		dirs = []string{cwd}
	}

	var found []string
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", dir, err)
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			repoPath := filepath.Join(dir, e.Name())
			if isGitRepo(repoPath) {
				found = append(found, repoPath)
			}
		}
		// also check if dir itself is a git repo
		if isGitRepo(dir) {
			found = append(found, dir)
		}
	}

	if len(found) == 0 {
		fmt.Println("no git repos found")
		return
	}

	results := make([]repoInfo, len(found))
	for i, path := range found {
		results[i] = inspectRepo(path, !reposNoFetch)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].name < results[j].name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPO\tBRANCH\tAHEAD\tBEHIND\tDIRTY\tSTATUS")

	for _, r := range results {
		if r.err != "" {
			fmt.Fprintf(w, "%s\t%s\t-\t-\t-\t✗ %s\n", r.name, r.branch, r.err)
			continue
		}

		dirty := ""
		if r.dirty {
			dirty = "✎"
		}

		status := "✓ clean"
		switch {
		case r.ahead > 0 && r.behind > 0:
			status = "⚡ diverged"
		case r.ahead > 0:
			status = "↑ push needed"
		case r.behind > 0:
			status = "↓ pull needed"
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n",
			r.name, r.branch, r.ahead, r.behind, dirty, status)
	}
	w.Flush()
}

func isGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func inspectRepo(path string, fetch bool) repoInfo {
	name := filepath.Base(path)
	info := repoInfo{path: path, name: name}

	branch, err := gitOutput(path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		info.err = "no HEAD"
		return info
	}
	info.branch = branch

	if fetch {
		gitOutput(path, "fetch", "--quiet") //nolint:errcheck
	}

	upstream, err := gitOutput(path, "rev-parse", "--abbrev-ref", "@{u}")
	if err != nil {
		info.err = "no upstream"
		return info
	}

	aheadStr, err := gitOutput(path, "rev-list", "--count", upstream+"..HEAD")
	if err == nil {
		info.ahead, _ = strconv.Atoi(aheadStr)
	}

	behindStr, err := gitOutput(path, "rev-list", "--count", "HEAD.."+upstream)
	if err == nil {
		info.behind, _ = strconv.Atoi(behindStr)
	}

	status, err := gitOutput(path, "status", "--porcelain")
	if err == nil {
		info.dirty = strings.TrimSpace(status) != ""
	}

	return info
}

func runReposPull(cmd *cobra.Command, args []string) {
	dirs := repoDirs
	if len(dirs) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot get cwd:", err)
			os.Exit(1)
		}
		dirs = []string{cwd}
	}

	var found []string
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			repoPath := filepath.Join(dir, e.Name())
			if isGitRepo(repoPath) {
				found = append(found, repoPath)
			}
		}
		if isGitRepo(dir) {
			found = append(found, dir)
		}
	}

	sort.Slice(found, func(i, j int) bool {
		return filepath.Base(found[i]) < filepath.Base(found[j])
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPO\tRESULT")

	for _, path := range found {
		name := filepath.Base(path)
		out, err := gitOutput(path, "pull", "--ff-only", "--quiet")
		if err != nil {
			// check if diverged vs other error
			if _, e2 := gitOutput(path, "rev-parse", "@{u}"); e2 != nil {
				fmt.Fprintf(w, "%s\t✗ no upstream\n", name)
			} else {
				fmt.Fprintf(w, "%s\t⚡ diverged — skipped\n", name)
			}
			continue
		}
		if out == "" || out == "Already up to date." {
			fmt.Fprintf(w, "%s\t✓ up to date\n", name)
		} else {
			fmt.Fprintf(w, "%s\t↓ pulled\n", name)
		}
	}
	w.Flush()
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
