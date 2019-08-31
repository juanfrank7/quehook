package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/forstmeier/comana/handlers"
	"github.com/forstmeier/comana/storage"
)

// HANDLER allows for build-time starter configuration
var HANDLER string

func starter(req handlers.Request) (events.APIGatewayProxyResponse, error) {
	s := storage.New()

	switch HANDLER {
	case "SAVE":
		return handlers.SaveData(req, s)
	case "LOAD":
		return handlers.LoadData(s)
	case "BACKFILL":
		i := handlers.NewInvoke()
		return handlers.BackfillData(req, i)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      500,
		Body:            "requested lambda type not available",
		IsBase64Encoded: false,
	}, errors.New("requested lambda type not available")
}

func main() {
	lambda.Start(starter)
}
