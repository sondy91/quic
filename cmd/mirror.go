package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var repos = []string{
	"Go-Projects", "VueECommerceProject", "ZeroCoolRepo", "ZeroCoolWebService",
	"advent-of-code", "astrowind-portfolio", "dotfiles", "duckie-boi",
	"dumpster-fire-driven-development", "making-a-dumpster-fire", "moonsweat-mountain",
	"personal-obsidian", "pluckMany", "quic", "scratchpad", "the-metrix", "vestramaximus",
}

type repoStatus struct {
	name       string
	githubSHA  string
	gitlabSHA  string
	githubErr  error
	gitlabErr  error
}

func init() {
	rootCmd.AddCommand(mirrorCmd)
	mirrorCmd.AddCommand(mirrorStatusCmd)
}

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Mirror sync utilities",
}

var mirrorStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitHub↔GitLab sync status for all repos",
	Run:   runMirrorStatus,
}

func runMirrorStatus(cmd *cobra.Command, args []string) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if githubToken == "" || gitlabToken == "" {
		fmt.Fprintln(os.Stderr, "error: GITHUB_TOKEN and GITLAB_TOKEN env vars required")
		os.Exit(1)
	}

	results := make([]repoStatus, len(repos))
	for i, r := range repos {
		results[i].name = r
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, repo := range repos {
		wg.Add(2)
		go func(idx int, r string) {
			defer wg.Done()
			sha, err := fetchGitHubSHA(r, githubToken)
			mu.Lock()
			results[idx].githubSHA = sha
			results[idx].githubErr = err
			mu.Unlock()
		}(i, repo)
		go func(idx int, r string) {
			defer wg.Done()
			sha, err := fetchGitLabSHA(r, gitlabToken)
			mu.Lock()
			results[idx].gitlabSHA = sha
			results[idx].gitlabErr = err
			mu.Unlock()
		}(i, repo)
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].name < results[j].name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPO\tGITHUB\tGITLAB\tSTATUS")

	drifted := 0
	for _, r := range results {
		github := shortSHA(r.githubSHA)
		gitlab := shortSHA(r.gitlabSHA)
		status := "✓ in sync"

		if r.githubErr != nil {
			github = "error"
			status = "✗ error"
			drifted++
		} else if r.gitlabErr != nil {
			gitlab = "error"
			status = "✗ error"
			drifted++
		} else if r.githubSHA != r.gitlabSHA {
			status = "✗ drift"
			drifted++
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.name, github, gitlab, status)
	}
	w.Flush()

	if drifted > 0 {
		fmt.Fprintf(os.Stderr, "\n%d repo(s) out of sync\n", drifted)
		os.Exit(1)
	}
}

func fetchGitHubSHA(repo, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/sondy91/%s/commits?per_page=1", repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var commits []struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return "", err
	}
	if len(commits) == 0 {
		return "", fmt.Errorf("no commits")
	}
	return commits[0].SHA, nil
}

func fetchGitLabSHA(repo, token string) (string, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/asonderman%%2F%s/repository/commits?per_page=1", repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var commits []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return "", err
	}
	if len(commits) == 0 {
		return "", fmt.Errorf("no commits")
	}
	return commits[0].ID, nil
}

func shortSHA(sha string) string {
	if len(sha) >= 7 {
		return sha[:7]
	}
	return sha
}
