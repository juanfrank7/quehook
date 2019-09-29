package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/forstmeier/quehook/query"
	"github.com/forstmeier/quehook/storage"
	"github.com/forstmeier/quehook/subscription"
	"github.com/forstmeier/quehook/table"
)

// HANDLER allows for build-time starter configuration
var HANDLER string

func starter(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s := storage.New()
	t := table.New()

	switch HANDLER {
	case "CREATE":
		return query.Create(req, t, s)
	case "RUN":
		return query.Run(s, t)
	case "DELETE":
		return query.Delete(req, t, s)
	case "SUBSCRIBE":
		return subscription.Subscribe(req, t)
	case "UNSUBSCRIBE":
		return subscription.Unsubscribe(req, t)
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
