package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/forstmeier/watchmyrepo/database"
)

type stats struct {
	Repos      int
	Datapoints int
	Watchers   int
}

type info struct {
	Stats stats
	Name  string
}

type errorMsg struct {
	Status int
	Error  string
}

// Display renders the application landing page
func Display(path string, database database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos, err := database.LoadRepos()
		if err != nil {
			tmpl := template.Must(template.ParseFiles(path + "error.html"))
			e := errorMsg{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			tmpl.Execute(w, e)
			return
		}

		s := info{
			Stats: stats{
				Repos: len(repos),
			},
		}

		tmpl := template.Must(template.ParseFiles(path + "index.html"))
		tmpl.Execute(w, s)
	}
}

func validateURL(r *http.Request) (string, error) {
	r.ParseForm()

	repo := r.Form["repo"][0]

	_, err := url.ParseRequestURI(repo)
	if err != nil {
		return "", err
	}

	if !strings.Contains(repo, "github.com") {
		return "", errors.New("'github.com' not found in url")
	}
	repoSplit := strings.Split(repo, "/")

	return repoSplit[3] + "/" + repoSplit[4], nil
}

// Submit processes new repo submissions to watch
func Submit(path string, database database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, err := validateURL(r)
		if err != nil {
			tmpl := template.Must(template.ParseFiles(path + "error.html"))
			e := errorMsg{
				Error:  err.Error(),
				Status: http.StatusBadRequest,
			}
			tmpl.Execute(w, e)
			return
		}

		if err := database.SaveRepo(repo); err != nil {
			tmpl := template.Must(template.ParseFiles(path + "error.html"))
			e := errorMsg{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			tmpl.Execute(w, e)
			return
		}

		repos, err := database.LoadRepos()
		if err != nil {
			tmpl := template.Must(template.ParseFiles(path + "error.html"))
			e := errorMsg{
				Error:  err.Error(),
				Status: http.StatusInternalServerError,
			}
			tmpl.Execute(w, e)
			return
		}

		resp := info{
			Name: repo,
			Stats: stats{
				Repos: len(repos),
			},
		}

		tmpl := template.Must(template.ParseFiles(path + "index.html"))
		tmpl.Execute(w, resp)
	}
}

// Data handles requests for watched repo stats
func Data(database database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := database.LoadStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		output, err := json.Marshal(stats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	}
}
