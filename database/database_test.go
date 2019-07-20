package database

import (
	"log"
	"os"
	"testing"

	"github.com/boltdb/bolt"

	"github.com/forstmeier/watchmyrepo/helpers"
)

func TestNewDB(t *testing.T) {
	boltDB, err := bolt.Open("test_newdb.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer boltDB.Close()

	db, err := New(boltDB)
	if err != nil {
		t.Errorf("description: error creating database, received: %s", err.Error())
	}

	if db == nil {
		t.Errorf("description: database created incorrectly, received: %+v", db)
	}

	if err := os.Remove("test_newdb.db"); err != nil {
		log.Fatal(err)
	}
}

func Test_repos(t *testing.T) {
	db, err := bolt.Open("test_repos.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := repoStatsDB{
		db: db,
	}

	repoName := "test-repo"
	if err := r.SaveRepo(repoName); err != nil {
		t.Errorf("description: error saving repo data, error: %v", err)
	}

	repos, err := r.LoadRepos()
	if err != nil {
		t.Errorf("description: error loading repo data, error: %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("description: incorrect repo count returned, received: %v, expected: %v", repos, []string{repoName})
	}

	if err := os.Remove("test_repos.db"); err != nil {
		log.Fatal(err)
	}
}

func Test_stats(t *testing.T) {
	db, err := bolt.Open("test_stats.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := repoStatsDB{
		db: db,
	}

	statsList := []helpers.Stat{
		{
			Name: "test-stat",
		},
	}

	if err := s.SaveStats(statsList); err != nil {
		t.Errorf("description: error saving stat data, error: %v", err)
	}

	stats, err := s.LoadStats()
	if err != nil {
		t.Errorf("description: error loading stat data, error: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("description: incorrect stat count returned, received: %+v, expected: %v", stats, statsList)
	}

	if err := os.Remove("test_stats.db"); err != nil {
		log.Fatal(err)
	}
}
