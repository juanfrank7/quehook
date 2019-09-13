package query

import (
        "github.com/aws/aws-lambda-go/events"
)

func Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

        // outline:
        // [ ] parse file from request
        // [ ] parse query name from request
        // [ ] check if query exists in dynamodb
        // - [ ] true:
        // - - [ ] return error
        // - [ ] false:
        // - - [ ] add query name to dynamodb
        // - - [ ] persist query file to s3
        // [ ] return success

}

func Run() (events.APIGatewayProxyResponse, error) {

        // outline:
        // [ ] read in all queries from s3
        // [ ] loop over queries
        // - [ ] read in subscribers to query from dynamodb
        // - [ ] run query on bq
        // - [ ] parse results to response json
        // - [ ] post results to webhook
        // [ ] return success

}

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

}
