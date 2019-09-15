package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/api/iterator"

	"github.com/forstmeier/comana/storage"
	"github.com/forstmeier/comana/table"
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
	return bigquery.NewClient(context.Background(), "comana")
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
func Delete(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// outline:
	// [ ] parse secret from request
	// [ ] parse query name from request
	// [ ] check if secret is valid
	// - [ ] true: continue
	// - [ ] false: return error
	// [ ] check if query exists in dynamodb
	// - [ ] true:
	// - - [ ] delete file from s3
	// - - [ ] delete query name from dynamodb
	// - [ ] false:
	// - - [ ] return error
	// [ ] return success

	return events.APIGatewayProxyResponse{}, nil
}
