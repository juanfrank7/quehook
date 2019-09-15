package query

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/forstmeier/comana/storage"
	"github.com/forstmeier/comana/table"
)

// Create adds a query to S3 for periodic execution
func Create(request events.APIGatewayProxyRequest, table table.Table, storage storage.Storage) (events.APIGatewayProxyResponse, error) {
	queryName := request.QueryStringParameters["query_name"]

	check, err := table.Get(queryName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error getting query table: " + err.Error(),
			IsBase64Encoded: false,
		}, fmt.Errorf("error getting query table: %s", err.Error())
	}

	if check == false {
		if err := table.Add("queries", queryName); err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error creating query: " + err.Error(),
				IsBase64Encoded: false,
			}, fmt.Errorf("error creating query: %s", err.Error())
		}

		if err := storage.PutFile(queryName, strings.NewReader(request.Body)); err != nil {
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

// Run executes all stored queries and returns results to subscribers
func Run() (events.APIGatewayProxyResponse, error) {

	// outline:
	// [ ] read in all queries from s3
	// [ ] loop over queries
	// - [ ] read in subscribers to query from dynamodb
	// - [ ] run query on bq
	// - [ ] parse results to response json
	// - [ ] post results to webhook
	// [ ] return success

	return events.APIGatewayProxyResponse{}, nil
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
