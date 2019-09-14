package table

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoDBClient interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

// Table provides helper methods for persisting/retrieving/deleting items
type Table interface {
	Add(table string, items ...string) error
	Get(table string, items ...string) (bool, error)
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
	input := &dynamodb.PutItemInput{}

	if table == "subscribers" {
		input = &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
				"subname": {
					S: aws.String(items[1]),
				},
				"target": {
					S: aws.String(items[2]),
				},
			},
			TableName: aws.String(table),
		}
	} else if table == "queries" {
		input = &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
			},
			TableName: aws.String(table),
		}
	}

	_, err := c.dynamodb.PutItem(input)
	if err != nil {
		return fmt.Errorf("put item error: %s", err.Error())
	}
	return nil
}

// Get retrieves an item from DynamoDB
func (c *Client) Get(table string, items ...string) (bool, error) {
	input := &dynamodb.GetItemInput{}

	if table == "subscribers" {
		input = &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
				"name": {
					S: aws.String(items[1]),
				},
			},
			TableName: aws.String(table),
		}
	} else if table == "queries" {
		input = &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
			},
			TableName: aws.String(table),
		}
	}

	_, err := c.dynamodb.GetItem(input)
	if err != nil {
		return false, fmt.Errorf("get item error: %s", err.Error())
	}
	return true, nil
}

// Remove deletes an item from DynamoDB
func (c *Client) Remove(table string, items ...string) error {
	input := &dynamodb.DeleteItemInput{}

	if table == "subscribers" {
		input = &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
				"name": {
					S: aws.String(items[1]),
				},
			},
			TableName: aws.String(table),
		}
	} else if table == "queries" {
		input = &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"query": {
					S: aws.String(items[0]),
				},
			},
			TableName: aws.String(table),
		}
	}

	_, err := c.dynamodb.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("delete item error: %s", err.Error())
	}
	return nil
}
