package table

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoDBClient interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error)
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

// Table provides helper methods for persisting/retrieving/deleting items
type Table interface {
	Add(table string, items ...string) error
	Get(table string, key string) ([]string, error)
	Remove(table string, item ...string) error
}

// Client implements the Table interface
type Client struct {
	dynamodb dynamoDBClient
}

// New generates a Table implementation with an active client
func New() Table {
	return &Client{
		dynamodb: dynamodb.New(session.New()),
	}
}

// Add puts a new item into DynamoDB
func (c *Client) Add(table string, items ...string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"query_name": {
				S: aws.String(items[0]),
			},
		},
	}

	if table == "subscribers" {
		input.Item["subscriber_email"] = &dynamodb.AttributeValue{
			S: aws.String(items[1]),
		}
		input.Item["subscriber_name"] = &dynamodb.AttributeValue{
			S: aws.String(items[2]),
		}
		input.Item["subscriber_target"] = &dynamodb.AttributeValue{
			S: aws.String(items[3]),
		}
	} else if table == "queries" {
		input.Item["author_name"] = &dynamodb.AttributeValue{
			S: aws.String(items[1]),
		}
		input.Item["author_email"] = &dynamodb.AttributeValue{
			S: aws.String(items[2]),
		}
	}

	_, err := c.dynamodb.PutItem(input)
	if err != nil {
		return fmt.Errorf("put item error: %s", err.Error())
	}
	return nil
}

// Get retrieves an item from DynamoDB
func (c *Client) Get(table string, key string) ([]string, error) {
	results := []string{}

	requestItems := map[string]*dynamodb.KeysAndAttributes{
		table: &dynamodb.KeysAndAttributes{
			ConsistentRead: aws.Bool(true),
			Keys: []map[string]*dynamodb.AttributeValue{
				map[string]*dynamodb.AttributeValue{
					"query_name": {
						S: aws.String(key),
					},
				},
			},
		},
	}

	for {
		input := &dynamodb.BatchGetItemInput{
			RequestItems: requestItems,
		}

		output, err := c.dynamodb.BatchGetItem(input)
		if err != nil {
			return nil, fmt.Errorf("get item error: %s", err.Error())
		}

		attributeName := ""
		if table == "subscribers" {
			attributeName = "subscriber_target"
		} else if table == "queries" {
			attributeName = "query_name"
		}

		for _, result := range output.Responses[table] {
			results = append(results, result[attributeName].GoString())
		}

		if output.UnprocessedKeys == nil {
			break
		}

		requestItems = output.UnprocessedKeys
	}

	return results, nil
}

// Remove deletes an item from DynamoDB
func (c *Client) Remove(table string, items ...string) error {
	// input := &dynamodb.DeleteItemInput{}
	//
	// if table == "subscribers" {
	// 	input = &dynamodb.DeleteItemInput{
	// 		Key: map[string]*dynamodb.AttributeValue{
	// 			"query": {
	// 				S: aws.String(items[0]),
	// 			},
	// 			"subname": {
	// 				S: aws.String(items[1]),
	// 			},
	// 		},
	// 		TableName: aws.String(table),
	// 	}
	// } else if table == "queries" {
	// 	input = &dynamodb.DeleteItemInput{
	// 		Key: map[string]*dynamodb.AttributeValue{
	// 			"query": {
	// 				S: aws.String(items[0]),
	// 			},
	// 		},
	// 		TableName: aws.String(table),
	// 	}
	// }
	//
	// _, err := c.dynamodb.DeleteItem(input)
	// if err != nil {
	// 	return fmt.Errorf("delete item error: %s", err.Error())
	// }
	return nil
}
