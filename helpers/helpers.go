package helpers

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

type helper interface {
	getStars(string, string) (int, error)
	getWatchers(string, string) (int, error)
	getCollaborators(string, string) (int, error)
	getDownloads(string, string) (int, error)
	getHealth(string, string) (int, error)
	getForks(string, string) (int, error)
	getContributors(string, string) (int, error)
}

// Stat holds repo-level data produced by helper collection functions
type Stat struct {
	Name          string
	Time          time.Time
	Stars         int
	Watchers      int
	Collaborators int
	Downloads     int
	Health        int
	Forks         int
	Contributors  int
}

// New creates an instance of the Help struct for external use
func New(client *github.Client) *Help {
	return &Help{
		client: client,
	}
}

// Help implements the helper interface methods
type Help struct {
	client *github.Client
}

func (h *Help) getStars(owner, repo string) (int, error) {
	count := 0
	ctx := context.Background()
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		stargazers, resp, err := h.client.Activity.ListStargazers(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}
		count += len(stargazers)

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	return count, nil
}

func (h *Help) getWatchers(owner, repo string) (int, error) {
	count := 0
	ctx := context.Background()
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		watchers, resp, err := h.client.Activity.ListWatchers(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}
		count += len(watchers)

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	return count, nil
}

func (h *Help) getCollaborators(owner, repo string) (int, error) {
	count := 0
	ctx := context.Background()
	opts := &github.ListCollaboratorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		collaborators, resp, err := h.client.Repositories.ListCollaborators(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}
		count += len(collaborators)

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	return count, nil
}

func (h *Help) getDownloads(owner, repo string) (int, error) {
	count := 0
	ctx := context.Background()
	opts := &github.ListOptions{
		PerPage: 100,
	}

	releases := []*github.RepositoryRelease{}
	resp := &github.Response{}
	var err error
	for {
		releases, resp, err = h.client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	assets := []github.ReleaseAsset{}
	for _, release := range releases {
		assets = release.Assets
	}

	for _, asset := range assets {
		count += *asset.DownloadCount
	}

	return count, nil
}

func (h *Help) getHealth(owner, repo string) (int, error) {
	ctx := context.Background()

	metrics, _, err := h.client.Repositories.GetCommunityHealthMetrics(ctx, owner, repo)
	if err != nil {
		return 0, err
	}

	return *metrics.HealthPercentage, nil
}

func (h *Help) getForks(owner, repo string) (int, error) {
	count := 0
	ctx := context.Background()
	opts := &github.RepositoryListForksOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		forks, resp, err := h.client.Repositories.ListForks(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}
		count += len(forks)

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	return count, nil
}

func (h *Help) getContributors(owner, repo string) (int, error) {
	ctx := context.Background()

	contributors, _, err := h.client.Repositories.ListContributorsStats(ctx, owner, repo)
	if err != nil {
		return 0, err
	}

	return len(contributors), nil
}

// GetDatapoints calls all helper functions and collects their output
func GetDatapoints(repo string, help helper) (*Stat, error) {
	owner := strings.Split(repo, "/")[0]
	name := strings.Split(repo, "/")[1]

	stars, err := help.getStars(owner, name)
	if err != nil {
		return nil, err
	}

	watchers, err := help.getWatchers(owner, name)
	if err != nil {
		return nil, err
	}

	collaborators, err := help.getCollaborators(owner, name)
	if err != nil {
		return nil, err
	}

	downloads, err := help.getDownloads(owner, name)
	if err != nil {
		return nil, err
	}

	health, err := help.getHealth(owner, name)
	if err != nil {
		return nil, err
	}

	forks, err := help.getForks(owner, name)
	if err != nil {
		return nil, err
	}

	contributors, err := help.getContributors(owner, name)
	if err != nil {
		return nil, err
	}

	s := &Stat{
		Name:          repo,
		Time:          time.Now(),
		Stars:         stars,
		Watchers:      watchers,
		Collaborators: collaborators,
		Downloads:     downloads,
		Health:        health,
		Forks:         forks,
		Contributors:  contributors,
	}

	return s, nil
}
