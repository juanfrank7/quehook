package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/forstmeier/comana/handlers"
	"github.com/forstmeier/comana/storage"
)

func router(req handlers.Request) (events.APIGatewayProxyResponse, error) {
	s := storage.New()

	switch req.HTTPMethod {
	case "GET":
		return handlers.LoadData(s)
	case "PUT":
		return handlers.BackfillData(req, s)
	default:
		return handlers.SaveData(req, s)
	}
}

func main() {
	lambda.Start(router)
}
