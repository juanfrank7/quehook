package main

import (
	"testing"

	"github.com/google/go-github/github"

	"github.com/forstmeier/watchmyrepo/helpers"
)

type dbMock struct{}

func (d *dbMock) SaveRepo(repo string) error {
	return nil
}

func (d *dbMock) LoadRepos() ([]string, error) {
	return []string{}, nil
}

func (d *dbMock) SaveStats(stats []helpers.Stat) error {
	return nil
}

func (d *dbMock) LoadStats() ([]map[string][]helpers.Stat, error) {
	return []map[string][]helpers.Stat{}, nil
}

func Test_serverFuncs(t *testing.T) {
	c := github.NewClient(nil)
	d := &dbMock{}
	s := newServer("secret", "127.0.0.1", c, d)

	if s.secret != "secret" && s.client != nil {
		t.Errorf("description: server created incorrectly")
	}

	s.start()
	s.stop()
}
