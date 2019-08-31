package handlers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/forstmeier/comana/storage"
)

// LoadData retrieves and returns GitHub Archive reports
func LoadData(s storage.Storage) (events.APIGatewayProxyResponse, error) {
	log.Println("load request")

	paths, err := s.GetPaths()
	log.Printf("paths count: %d", len(paths))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error loading report filepaths: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	pathsObject := map[string][]string{
		"paths": paths,
	}

	output, err := json.Marshal(pathsObject)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error marshalling output: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	log.Println("load successful")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"result_count": strconv.Itoa(len(paths)),
		},
		Body:            string(output),
		IsBase64Encoded: false,
	}, err
}
