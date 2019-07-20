package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

func TestNew(t *testing.T) {
	c := github.NewClient(nil)
	h := New(c)

	if h.client == nil {
		t.Errorf("description: new help function returning incorrect values")
	}
}

func Test_helpers(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/tatooine/homestead/stargazers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"user":{"login":"skywalker"}}]`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/subscribers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"user":{"login":"kenobi"}}]`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/collaborators", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"user":{"login":"lars"}},{"user":{"login":"whitesun"}}]`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/releases", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"assets":[{"download_count":2}]}]`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/community/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"health_percentage":0}`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/forks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"tosche-station"}]`)
	})

	mux.HandleFunc("/repos/tatooine/homestead/stats/contributors", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"author":{"login":"r2"}},{"author":{"login":"3po"}}]`)
	})

	server := httptest.NewServer(mux)

	c := github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	c.BaseURL = url
	c.UploadURL = url

	h := Help{
		client: c,
	}

	t.Run("test get stargazers", func(t *testing.T) {
		count, err := h.getStars("tatooine", "homestead")
		if count != 1 {
			t.Errorf("bulk stargazers retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk stargazers retrieval error: %s", err.Error())
		}
	})

	t.Run("test get watchers", func(t *testing.T) {
		count, err := h.getWatchers("tatooine", "homestead")
		if count != 1 {
			t.Errorf("bulk watchers retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk watchers retrieval error: %s", err.Error())
		}
	})

	t.Run("test get collaborators", func(t *testing.T) {
		count, err := h.getCollaborators("tatooine", "homestead")
		if count != 2 {
			t.Errorf("bulk collaborators retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk collaborators retrieval error: %s", err.Error())
		}
	})

	t.Run("test get downloads", func(t *testing.T) {
		count, err := h.getDownloads("tatooine", "homestead")
		if count != 2 {
			t.Errorf("bulk downloads retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk downloads retrieval error: %s", err.Error())
		}
	})

	t.Run("test get health", func(t *testing.T) {
		percent, err := h.getHealth("tatooine", "homestead")
		if percent != 0 {
			t.Errorf("bulk health retrieval incorrect percent: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk health retrieval error: %s", err.Error())
		}
	})

	t.Run("test get forks", func(t *testing.T) {
		count, err := h.getForks("tatooine", "homestead")
		if count != 1 {
			t.Errorf("bulk forks retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk forks retrieval error: %s", err.Error())
		}
	})

	t.Run("test get contributors", func(t *testing.T) {
		count, err := h.getContributors("tatooine", "homestead")
		if count != 2 {
			t.Errorf("bulk contributors retrieval miscount: %s", err.Error())
		}

		if err != nil {
			t.Errorf("bulk contributors retrieval error: %s", err.Error())
		}
	})
}

type helpMock struct {
	count int
	err   error
}

func (h *helpMock) getStars(o, r string) (int, error) {
	if r != "naboo" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getWatchers(o, r string) (int, error) {
	if r != "coruscant" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getCollaborators(o, r string) (int, error) {
	if r != "mustafar" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getDownloads(o, r string) (int, error) {
	if r != "tatooine" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getHealth(o, r string) (int, error) {
	if r != "hoth" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getForks(o, r string) (int, error) {
	if r != "endor" {
		return 1, nil
	}
	return h.count, h.err
}

func (h *helpMock) getContributors(o, r string) (int, error) {
	if r != "geonosis" {
		return 1, nil
	}
	return h.count, h.err
}

func Test_callHelpers(t *testing.T) {
	tests := []struct {
		desc   string
		repo   string
		count  int
		err    error
		output *Stat
	}{
		{
			"get stars error scenario",
			"mid-rim/naboo",
			0,
			errors.New("error getting stars"),
			nil,
		},
		{
			"get watchers error scenario",
			"galactic-core/coruscant",
			0,
			errors.New("error getting watchers"),
			nil,
		},
		{
			"get collaborators error scenario",
			"outer-rim/mustafar",
			0,
			errors.New("error getting collaborators"),
			nil,
		},
		{
			"get downloads error scenario",
			"outer-rim/tatooine",
			0,
			errors.New("error getting downloads"),
			nil,
		},
		{
			"get health error scenario",
			"outer-rim/hoth",
			0,
			errors.New("error getting health"),
			nil,
		},
		{
			"get forks error scenario",
			"outer-rim/endor",
			0,
			errors.New("error getting forks"),
			nil,
		},
		{
			"get contributors error scenario",
			"outer-rim/geonosis",
			0,
			errors.New("error getting contributors"),
			nil,
		},
		{
			"successful stat retrieval scenario",
			"wild-space/kamino",
			0,
			nil,
			&Stat{
				Stars:         1,
				Watchers:      1,
				Collaborators: 1,
				Downloads:     1,
				Health:        1,
				Forks:         1,
				Contributors:  1,
			},
		},
	}

	for _, test := range tests {
		output, err := GetDatapoints(test.repo, &helpMock{
			count: test.count,
			err:   test.err,
		})

		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("description: %s, incorrect error message, received: %s, expected: %s", test.desc, err.Error(), test.err.Error())
		}

		if test.output != nil && output != nil {
			if output.Stars != test.output.Stars || output.Downloads != test.output.Downloads || output.Forks != test.output.Forks {
				t.Errorf("description: %s, incorrect stats, received: %+v, expected: %+v", test.desc, output, test.output)
			}
		}
	}
}
