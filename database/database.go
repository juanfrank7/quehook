package database

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"

	"github.com/forstmeier/watchmyrepo/helpers"
)

// DB provides stat and repo data persistence functionality
type DB interface {
	SaveRepo(repo string) error
	LoadRepos() ([]string, error)
	SaveStats(stats []helpers.Stat) error
	LoadStats() ([]map[string][]helpers.Stat, error)
}

// New generates an implementation of the DB interface
func New(boltDB *bolt.DB) (DB, error) {
	if err := boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("repos"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("stats"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &repoStatsDB{
		db: boltDB,
	}, nil
}

type repoStatsDB struct {
	db *bolt.DB
}

func (r *repoStatsDB) SaveRepo(repo string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("repos"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		currentTime := time.Now().String()
		if err := bucket.Put([]byte(repo), []byte(currentTime)); err != nil {
			return fmt.Errorf("failed to insert '%s': %v", repo, err)
		}
		return nil
	})
}

func (r *repoStatsDB) LoadRepos() ([]string, error) {
	var repos []string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("repos"))
		b.ForEach(func(k, v []byte) error {
			repos = append(repos, string(k))
			return nil
		})
		return nil
	})

	return repos, err
}

func (r *repoStatsDB) SaveStats(stats []helpers.Stat) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("stats"))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		statBytes, err := json.Marshal(stats)
		if err != nil {
			return err
		}

		key := time.Now().String()
		if err := bucket.Put([]byte(key), []byte(statBytes)); err != nil {
			return fmt.Errorf("failed to insert '%s': %v", key, err)
		}

		return nil
	})
}

func (r *repoStatsDB) LoadStats() ([]map[string][]helpers.Stat, error) {
	var stats []map[string][]helpers.Stat

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("stats"))
		b.ForEach(func(k, v []byte) error {
			s := []helpers.Stat{}
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}

			repoStat := map[string][]helpers.Stat{
				string(k): s,
			}

			stats = append(stats, repoStat)
			return nil
		})
		return nil
	})

	return stats, err
}

func getDatapoints(database DB, client *github.Client) error {
	repos, err := database.LoadRepos()
	if err != nil {
		return err
	}

	help := helpers.New(client)

	repoStats := []helpers.Stat{}

	var wg sync.WaitGroup
	wg.Add(len(repos))
	for _, repo := range repos {
		go func(val string) {
			s, err := helpers.GetDatapoints(val, help)
			_ = err
			repoStats = append(repoStats, *s)
		}(repo)
	}
	wg.Wait()

	if err := database.SaveStats(repoStats); err != nil {
		return err
	}

	return nil
}
