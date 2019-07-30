package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/forstmeier/watchmyrepo/helpers"
)

type dbMock struct {
	stats   []map[string][]helpers.Stat
	repos   []string
	errLoad error
	errSave error
}

func (s *dbMock) SaveRepo(repo string) error {
	return s.errSave
}

func (s *dbMock) LoadRepos() ([]string, error) {
	return s.repos, s.errLoad
}

func (s *dbMock) SaveStats(stats []helpers.Stat) error {
	return s.errSave
}

func (s *dbMock) LoadStats() ([]map[string][]helpers.Stat, error) {
	return s.stats, s.errLoad
}

func Test_display(t *testing.T) {
	tests := []struct {
		desc   string
		repos  []string
		err    error
		status int
		output string
	}{
		{
			"error loading repos",
			nil,
			errors.New("error loading repos"),
			200,
			`<p>Circle back to the <b><a href="/">home page</a></b> to try again! Here's a cactus for your troubles. 🌵</p>`,
		},
		{
			"successful request",
			[]string{
				"nerf-herder/millenium-falcon",
			},
			nil,
			200,
			"<h1>Watch My Repo 🔭</h1>",
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()

		db := &dbMock{
			repos:   test.repos,
			errLoad: test.err,
		}

		handler := http.HandlerFunc(Display("../site/", db))
		handler.ServeHTTP(rec, req)

		if status := rec.Code; status != test.status {
			t.Errorf("description: %s, incorrect status code, received: %d, expected: %d", test.desc, status, test.status)
		}

		if received := rec.Body.String(); !strings.Contains(received, test.output) {
			t.Errorf("description: %s, returned incorrect response, received: %v, expected: %v", test.desc, received, test.output)
		}
	}

}

func Test_validateURL(t *testing.T) {
	tests := []struct {
		url  string
		repo string
		err  error
	}{
		{
			"test.com",
			"",
			errors.New("parse test.com: invalid URI for request"),
		},
		{
			"https://test.com",
			"",
			errors.New("'github.com' not found in url"),
		},
		{
			"https://github.com/test/test",
			"test/test",
			nil,
		},
	}

	for _, test := range tests {

		req, err := http.NewRequest("POST", "test-url", nil)
		if err != nil {
			log.Fatalf("error creating request: %s", err.Error())
		}
		form := url.Values{}
		form.Add("repo", test.url)
		req.PostForm = form

		output, err := validateURL(req)

		if output != test.repo {
			t.Errorf("description: incorrect repo value, received: %s, expected: %s", output, test.repo)
		}

		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("description: incorrect error value, received: %s, expected: %s", err.Error(), test.err.Error())
		}
	}
}

func Test_submit(t *testing.T) {
	tests := []struct {
		desc    string
		path    string
		repos   []string
		errSave error
		errLoad error
		output  string
	}{
		{
			"incorrect filepath",
			"/incorrect/",
			nil,
			errors.New("incorrect filepath"),
			nil,
			"<p><b>Status</b>: 500 <b>Error</b>: incorrect filepath</p>",
		},
		{
			"error saving repo",
			"../site/",
			nil,
			errors.New("error saving repo"),
			nil,
			"<p><b>Status</b>: 500 <b>Error</b>: error saving repo</p>",
		},
		{
			"error loading repo",
			"../site/",
			nil,
			nil,
			errors.New("error loading repo"),
			"<p><b>Status</b>: 500 <b>Error</b>: error loading repo</p>",
		},
		{
			"successful request",
			"../site/",
			[]string{
				"test-owner/test-repo",
			},
			nil,
			nil,
			`<div class="dynamic response">test-owner/test-repo submitted successfully</div>`,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		form := url.Values{}
		form.Add("repo", "https://github.com/test-owner/test-repo")
		req.PostForm = form
		rec := httptest.NewRecorder()

		db := &dbMock{
			repos:   test.repos,
			errSave: test.errSave,
			errLoad: test.errLoad,
		}

		handler := http.HandlerFunc(Submit("../site/", db))
		handler.ServeHTTP(rec, req)

		if received := rec.Body.String(); !strings.Contains(received, test.output) {
			t.Errorf("description: %s, display handler returned incorrect response, received: %v, expected: %v", test.desc, received, test.output)
		}
	}
}

func Test_data(t *testing.T) {
	tests := []struct {
		desc   string
		stats  []map[string][]helpers.Stat
		err    error
		status int
		output string
	}{
		{
			"error loading stats",
			nil,
			errors.New("error loading stats"),
			500,
			"error loading stats",
		},
		{
			"successful request",
			[]map[string][]helpers.Stat{
				map[string][]helpers.Stat{
					"owner/repo": []helpers.Stat{
						helpers.Stat{
							Name: "owner/repo",
						},
					},
				},
			},
			nil,
			200,
			`[{"owner/repo":[{"Name":"owner/repo","Time":"0001-01-01T00:00:00Z","Stars":0,"Watchers":0,"Collaborators":0,"Downloads":0,"Health":0,"Forks":0,"Contributors":0}]}]`,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()

		db := &dbMock{
			stats:   test.stats,
			errLoad: test.err,
		}

		handler := http.HandlerFunc(Data(db))
		handler.ServeHTTP(rec, req)

		status := rec.Code
		if status != test.status {
			t.Errorf("description: %s, incorrect status returned, received: %d, expected: %d", test.desc, status, test.status)
		}

		if received := rec.Body.String(); !strings.Contains(received, test.output) {
			t.Errorf("description: %s, incorrect json output, received: %v, expected: %v", test.desc, received, test.output)
		}
	}
}
