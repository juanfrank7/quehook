package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/go-github/github"

	"github.com/forstmeier/watchmyrepo/helpers"
)

// DB provides stat and repo data persistence functionality
type DB interface {
	SaveRepo(ctx context.Context, repo string) error
	LoadRepos(ctx context.Context) ([]string, error)
	SaveStats(ctx context.Context, stats []helpers.Stat) error
	LoadStats(ctx context.Context) ([]map[string][]helpers.Stat, error)
	SaveWatcher(ctx context.Context, id string) error
	LoadWatchers(ctx context.Context) ([]string, error)
}

// New generates an implementation of the DB interface
func New() (DB, error) {
	return &db{}, nil
}

type db struct{}

func newClient(ctx context.Context) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, "watchmyrepo-248001")
	if err != nil {
		return nil, fmt.Errorf("error creating datastore client: %s", err.Error())
	}
	return client, nil
}

// Datapoint is use internally to process Datastore interactions
type Datapoint struct {
	Repo    string
	Watcher string
	Stats   []helpers.Stat
}

func (d *db) SaveRepo(ctx context.Context, repo string) error {
	key := datastore.IncompleteKey("repo", nil)
	if key == nil {
		return errors.New("error creating repo key")
	}

	dp := Datapoint{
		Repo: repo,
	}

	client, err := newClient(ctx)
	if err != nil {
		return err
	}

	if _, err := client.Put(ctx, key, &dp); err != nil {
		return fmt.Errorf("error saving repo value: %s", err.Error())
	}

	return nil
}

func (d *db) LoadRepos(ctx context.Context) ([]string, error) {
	query := datastore.NewQuery("repo")
	datapoints := []Datapoint{}

	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	_, err = client.GetAll(ctx, query, &datapoints)
	if err != nil {
		return nil, fmt.Errorf("error loading repos: %s", err.Error())
	}

	repos := []string{}
	for _, datapoint := range datapoints {
		repos = append(repos, datapoint.Repo)
	}

	return repos, nil
}

func (d *db) SaveStats(ctx context.Context, stats []helpers.Stat) error {
	currentTime := time.Now().String()
	key := datastore.NameKey("stat", currentTime, nil)
	if key == nil {
		return errors.New("error creating stats key")
	}

	datapoint := Datapoint{
		Stats: stats,
	}

	client, err := newClient(ctx)
	if err != nil {
		return err
	}

	if _, err := client.Put(ctx, key, &datapoint); err != nil {
		return fmt.Errorf("error saving stats value: %s", err.Error())
	}

	return nil
}

func (d *db) LoadStats(ctx context.Context) ([]map[string][]helpers.Stat, error) {
	query := datastore.NewQuery("stat")

	output := []map[string][]helpers.Stat{}
	datapoints := []Datapoint{}

	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	keys, err := client.GetAll(ctx, query, &datapoints)
	if err != nil {
		return nil, fmt.Errorf("error loading stats: %s", err.Error())
	}
	for i, key := range keys {
		stat := map[string][]helpers.Stat{
			key.String(): datapoints[i].Stats,
		}
		output = append(output, stat)
	}

	return output, nil
}

func (d *db) SaveWatcher(ctx context.Context, watcher string) error {
	key := datastore.IncompleteKey("watcher", nil)
	if key == nil {
		return errors.New("error creating watcher key")
	}

	datapoint := Datapoint{
		Watcher: watcher,
	}

	client, err := newClient(ctx)
	if err != nil {
		return err
	}

	if _, err := client.Put(ctx, key, &datapoint); err != nil {
		return fmt.Errorf("error saving watcher value: %s", err.Error())
	}

	return nil
}

func (d *db) LoadWatchers(ctx context.Context) ([]string, error) {
	query := datastore.NewQuery("watcher")

	datapoints := []Datapoint{}

	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	_, err = client.GetAll(ctx, query, &datapoints)
	if err != nil {
		return nil, fmt.Errorf("error loading watcher: %s", err.Error())
	}

	watchers := []string{}
	for _, datapoint := range datapoints {
		watchers = append(watchers, datapoint.Repo)
	}

	return watchers, nil
}

func getDatapoints(ctx context.Context, database DB, client *github.Client) error {
	repos, err := database.LoadRepos(ctx)
	if err != nil {
		return err
	}

	helper := helpers.New(client)

	repoStats := []helpers.Stat{}

	var wg sync.WaitGroup
	wg.Add(len(repos))
	for _, repo := range repos {
		go func(val string) {
			s, err := helpers.GetDatapoints(val, helper)
			_ = err
			repoStats = append(repoStats, *s)
			wg.Done()
		}(repo)
	}
	wg.Wait()

	if err := database.SaveStats(ctx, repoStats); err != nil {
		return err
	}

	return nil
}
