package database

import (
	"testing"

	"github.com/forstmeier/watchmyrepo/helpers"
	"google.golang.org/appengine/aetest"
)

func TestNewDB(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Errorf("description: error creating database, received: %s", err.Error())
	}

	if db == nil {
		t.Errorf("description: database created incorrectly, received: %+v", db)
	}
}

func Test_repos(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	d := db{}

	repoName := "test-repo"
	if err := d.SaveRepo(ctx, repoName); err != nil {
		t.Errorf("description: error saving repo data, error: %v", err)
	}

	repos, err := d.LoadRepos(ctx)
	if err != nil {
		t.Errorf("description: error loading repo data, error: %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("description: incorrect repo count returned, received: %v, expected: %v", repos, []string{repoName})
	}
}

func Test_stats(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	d := db{}
	statsList := []helpers.Stat{
		{
			Name: "test-stat",
		},
	}

	if err := d.SaveStats(ctx, statsList); err != nil {
		t.Errorf("description: error saving stat data, error: %v", err)
	}

	stats, err := d.LoadStats(ctx)
	if err != nil {
		t.Errorf("description: error loading stat data, error: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("description: incorrect stat count returned, received: %+v, expected: %v", stats, statsList)
	}
}

func Test_watchers(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	d := db{}
	id := "test-id"
	if err := d.SaveWatcher(ctx, id); err != nil {
		t.Errorf("description: error saving watcher data, error: %v", err)
	}

	watchers, err := d.LoadWatchers(ctx)
	if err != nil {
		t.Errorf("description: error loading stat data, error: %v", err)
	}

	if len(watchers) != 1 {
		t.Errorf("description: incorrect stat count returned, received: %+v, expected: %v", id, watchers)
	}
}
