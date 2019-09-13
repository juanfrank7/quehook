package subscription

import (
        "github.com/aws/aws-lambda-go/events"
)

func Subscribe(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

        // outline:
        // [ ] parse query name from request
        // [ ] parse webhook target / email target values from request
        // [ ] check if target values are invalid
        // - [ ] true: return error
        // - [ ] false: continue
        // [ ] check if query exists in dynamodb
        // - [ ] true:
        // - - [ ] add name to subscription table in dynamodb
        // - - [ ] add target values to subscription table in dynamodb
        // - [ ] false: return error
        // [ ] return success
        
}

func Unsubscribe(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

        // outline:
        // [ ] parse query name from request
        // [ ] check if name exists in dynamodb subscription table
        // - [ ] true:
        // - - [ ] delete item from dynamodb subscription table
        // - [ ] false:
        // - - [ ] return error
        // [ ] return success

}
