package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/api/iterator"

	"github.com/forstmeier/quehook/storage"
	"github.com/forstmeier/quehook/table"
)

// Create adds a query to S3 for periodic execution
func Create(request events.APIGatewayProxyRequest, t table.Table, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	queryName := request.QueryStringParameters["query_name"]

	_, check, err := t.Get("queries", queryName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error getting query table: " + err.Error(),
			IsBase64Encoded: false,
		}, fmt.Errorf("error getting query table: %s", err.Error())
	}

	if check == false {
		if err := t.Add("queries", queryName); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error creating query: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error creating query: %s", err.Error())
		}

		if err := s.PutFile(queryName, strings.NewReader(request.Body)); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error putting query file: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error putting query file: %s", err.Error())
		}
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            "query already exists",
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
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
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error listing query files: " + err.Error(),
			IsBase64Encoded: false,
		}, fmt.Errorf("error listing query files: %s", err.Error())
	}

	for _, query := range queries {
		file, err := s.GetFile(query)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error getting query file: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error getting query file: %s", err.Error())
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)

		q := bq.Query(buf.String())

		rows := [][]bigquery.Value{}
		itr, err := q.Read(context.Background())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error reading query: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error reading query: %s", err.Error())
		}

		for {
			var row []bigquery.Value
			err := itr.Next(&row)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error iterating query results: " + err.Error(),
					IsBase64Encoded: false,
				}, fmt.Errorf("error iterating query results: %s", err.Error())
			}

			rows = append(rows, row)
		}

		output, err := json.Marshal(rows)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error marshalling output: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error marshalling output: %s", err.Error())
		}

		subscribers, _, err := t.Get("subscribers", query)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error getting subscribers: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error getting subscribers: %s", err.Error())
		}

		client := &http.Client{}
		for _, subscriber := range subscribers {
			req, err := http.NewRequest("POST", subscriber, bytes.NewBuffer(output))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error posting results: " + err.Error(),
					IsBase64Encoded: false,
				}, fmt.Errorf("error posting results: %s", err.Error())
			}
			_ = resp // TEMP
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}

// Delete removes a query from S3 - internal use only
func Delete(request events.APIGatewayProxyRequest, t table.Table, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	if request.Headers["QUEHOOK_SECRET"] != os.Getenv("QUEHOOK_SECRET") {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "incorrect secret received: " + request.Headers["QUEHOOK_SECRET"],
			IsBase64Encoded: false,
		}, fmt.Errorf("incorrect secret received: %s", request.Headers["QUEHOOK_SECRET"])
	}

	body := struct {
		query string
	}{}

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error parsing request body: " + err.Error(),
			IsBase64Encoded: false,
		}, fmt.Errorf("incorrect secret received: %s", err.Error())
	}

	_, check, err := t.Get("queries", body.query)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error getting query: " + err.Error(),
			IsBase64Encoded: false,
		}, fmt.Errorf("incorrect getting query: %s", err.Error())
	}

	if check {
		if err := s.DeleteFile(body.query); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error deleting query file: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("incorrect deleting query file: %s", err.Error())
		}

		if err := t.Remove("queries", body.query); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error removing query item: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("incorrect removing query item: %s", err.Error())
		}

		if err := t.Remove("subscribers", body.query, ""); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error removing subscribers items: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("incorrect removing subscribers items: %s", err.Error())
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}
