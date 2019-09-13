package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func starter() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(starter)
}
