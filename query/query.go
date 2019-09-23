package query

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/api/iterator"

	"github.com/forstmeier/quehook/storage"
	"github.com/forstmeier/quehook/table"
)

func createResponse(code int, msg string) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		StatusCode:      code,
		Body:            msg,
		IsBase64Encoded: false,
	}

	if msg == "success" || msg == "query already exists" {
		return resp, nil
	}

	return resp, errors.New(msg)
}

// Create adds a query to S3 for periodic execution
func Create(request events.APIGatewayProxyRequest, t table.Table, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	queryName := request.QueryStringParameters["query_name"]

	output, err := t.Get("queries", queryName)
	if err != nil {
		return createResponse(500, "error getting query table: "+err.Error())
	}

	if len(output) == 0 {
		if err := t.Add("queries", queryName); err != nil {
			return createResponse(500, "error creating query: "+err.Error())
		}

		if err := s.PutFile(queryName, strings.NewReader(request.Body)); err != nil {
			return createResponse(500, "error putting query file: "+err.Error())
		}
	} else {
		return createResponse(200, "query already exists")
	}

	return createResponse(200, "success")
}

// BQClient wraps BigQuery methods and functionality
type BQClient interface {
	Query(query string) *bigquery.Query
}

// NewClient creates a new BigQuery client implementation
func NewClient() (BQClient, error) {
	return bigquery.NewClient(context.Background(), "quehook")
}

// Run executes all stored queries and returns results to subscribers
func Run(bq BQClient, s storage.Storage, t table.Table) (events.APIGatewayProxyResponse, error) {
	queries, err := s.GetPaths()
	if err != nil {
		return createResponse(500, "error listing query files: "+err.Error())
	}

	for _, query := range queries {
		file, err := s.GetFile(query)

		if err != nil {
			return createResponse(500, "error getting query file: "+err.Error())
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)

		q := bq.Query(buf.String())

		rows := [][]bigquery.Value{}
		itr, err := q.Read(context.Background())
		if err != nil {
			return createResponse(500, "error reading query: "+err.Error())
		}

		for {
			var row []bigquery.Value
			err := itr.Next(&row)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return createResponse(500, "error iterating query results: "+err.Error())
			}

			rows = append(rows, row)
		}

		output, err := json.Marshal(rows)
		if err != nil {
			return createResponse(500, "error marshalling output: "+err.Error())
		}

		subscribers, err := t.Get("subscribers", query)
		if err != nil {
			return createResponse(500, "error getting subscribers: "+err.Error())
		}

		client := &http.Client{}
		for _, subscriber := range subscribers {
			req, err := http.NewRequest("POST", subscriber, bytes.NewBuffer(output))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return createResponse(500, "error posting results: "+err.Error())
			}
			_ = resp // TEMP
		}
	}

	return createResponse(200, "success")
}

// Delete removes a query from S3 - internal use only
func Delete(request events.APIGatewayProxyRequest, t table.Table, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	if request.Headers["QUEHOOK_SECRET"] != os.Getenv("QUEHOOK_SECRET") {
		return createResponse(500, "incorrect secret received: "+request.Headers["QUEHOOK_SECRET"])
	}

	body := struct {
		query string
	}{}

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return createResponse(500, "error parsing request body: "+err.Error())
	}

	output, err := t.Get("queries", body.query)

	if err != nil {
		return createResponse(500, "error getting query: "+err.Error())
	}

	if len(output) > 0 {
		if err := s.DeleteFile(body.query); err != nil {
			return createResponse(500, "error deleting query file: "+err.Error())
		}

		if err := t.Remove("queries", body.query); err != nil {
			return createResponse(500, "error removing query item: "+err.Error())
		}

		if err := t.Remove("subscribers", body.query); err != nil {
			return createResponse(500, "error removing subscribers items: "+err.Error())
		}
	}

	return createResponse(200, "success")
}
