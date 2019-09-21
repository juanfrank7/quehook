package table

import (
	"errors"
	"testing"

	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type tableMock struct {
	putItemOutput    *dynamodb.PutItemOutput
	putItemError     error
	getItemOutput    *dynamodb.BatchGetItemOutput
	getItemError     error
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemError  error
}

func (mock *tableMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return mock.putItemOutput, mock.putItemError
}

func (mock *tableMock) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	return mock.getItemOutput, mock.getItemError
}

func (mock *tableMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return mock.deleteItemOutput, mock.deleteItemError
}

func TestAdd(t *testing.T) {
	tests := []struct {
		desc          string
		table         string
		items         []string
		putItemOutput *dynamodb.PutItemOutput // kept for possible method expansion
		putItemError  error
		err           string
	}{
		{
			desc:  "put item error",
			table: "queries",
			items: []string{
				"newPowerConverters",
				"luke",
				"luke@lars.homestead",
			},
			putItemOutput: nil,
			putItemError:  errors.New("mock put error"),
			err:           "put item error: mock put error",
		},
		{
			desc:  "successful subscribers invocation",
			table: "subscribers",
			items: []string{
				"newChores",
				"luke@lars.homstead",
				"luke",
				"https://holonet.com/skywalker",
			},
			putItemOutput: nil,
			putItemError:  nil,
			err:           "",
		},
		{
			desc:  "successful queries invocation",
			table: "queries",
			items: []string{
				"lotsOfTrouble",
				"r2-d2",
				"blue-and-white@royalengineers.nb",
			},
			putItemOutput: nil,
			putItemError:  nil,
			err:           "",
		},
	}

	for _, test := range tests {
		c := &Client{
			dynamodb: &tableMock{
				putItemOutput: test.putItemOutput,
				putItemError:  test.putItemError,
			},
		}

		if err := c.Add(test.table, test.items...); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func TestGet(t *testing.T) {
	// tests := []struct {
	// 	desc          string
	// 	table         string
	// 	items         []string
	// 	getItemOutput *dynamodb.BatchGetItemOutput
	// 	getItemError  error
	// 	check         bool
	// 	err           string
	// }{
	// 	{
	// 		desc:  "get item error",
	// 		table: "queries",
	// 		items: []string{
	// 			"query",
	// 		},
	// 		getItemOutput: nil,
	// 		getItemError:  errors.New("mock get error"),
	// 		check:         false,
	// 		err:           "get item error: mock get error",
	// 	},
	// 	{
	// 		desc:  "successful subscribers invocation",
	// 		table: "subscribers",
	// 		items: []string{
	// 			"query",
	// 			"subname",
	// 			"target",
	// 		},
	// 		getItemOutput: &dynamodb.BatchGetItemOutput{
	// 			Responses: map[string][]map[string]*dynamodb.AttributeValue{
	// 				"subscribers": []map[string]*dynamodb.AttributeValue{
	// 					map[string]*dynamodb.AttributeValue{
	// 						"query": {
	// 							S: aws.String("test-query"),
	// 						},
	// 						"target": {
	// 							S: aws.String("test-target"),
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		getItemError: nil,
	// 		check:        true,
	// 		err:          "",
	// 	},
	// 	{
	// 		desc:  "successful queries invocation",
	// 		table: "queries",
	// 		items: []string{
	// 			"query",
	// 		},
	// 		getItemOutput: &dynamodb.BatchGetItemOutput{
	// 			Responses: map[string][]map[string]*dynamodb.AttributeValue{
	// 				"subscribers": []map[string]*dynamodb.AttributeValue{
	// 					map[string]*dynamodb.AttributeValue{
	// 						"query": {
	// 							S: aws.String("test-query"),
	// 						},
	// 						"target": {
	// 							S: aws.String("test-target"),
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		getItemError: nil,
	// 		check:        true,
	// 		err:          "",
	// 	},
	// }
	//
	// for _, test := range tests {
	// 	c := &Client{
	// 		dynamodb: &tableMock{
	// 			getItemOutput: test.getItemOutput,
	// 			getItemError:  test.getItemError,
	// 		},
	// 	}
	//
	// 	_, check, err := c.Get(test.table, test.items...)
	//
	// 	if check != test.check {
	// 		t.Errorf("description: %s, check received: %t, expected: %t", test.desc, check, test.check)
	// 	}
	//
	// 	if err != nil && err.Error() != test.err {
	// 		t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
	// 	}
	// }
}

func TestRemove(t *testing.T) {
	// tests := []struct {
	// 	desc             string
	// 	table            string
	// 	items            []string
	// 	deleteItemOutput *dynamodb.DeleteItemOutput // kept for possible method expansion
	// 	deleteItemError  error
	// 	err              string
	// }{
	// 	{
	// 		desc:  "delete item error",
	// 		table: "queries",
	// 		items: []string{
	// 			"query",
	// 		},
	// 		deleteItemOutput: nil,
	// 		deleteItemError:  errors.New("mock delete error"),
	// 		err:              "delete item error: mock delete error",
	// 	},
	// 	{
	// 		desc:  "successful subscribers invocation",
	// 		table: "subscribers",
	// 		items: []string{
	// 			"query",
	// 			"subname",
	// 			"target",
	// 		},
	// 		deleteItemOutput: nil,
	// 		deleteItemError:  nil,
	// 		err:              "",
	// 	},
	// 	{
	// 		desc:  "successful queries invocation",
	// 		table: "queries",
	// 		items: []string{
	// 			"query",
	// 		},
	// 		deleteItemOutput: nil,
	// 		deleteItemError:  nil,
	// 		err:              "",
	// 	},
	// }
	//
	// for _, test := range tests {
	// 	c := &Client{
	// 		dynamodb: &tableMock{
	// 			deleteItemOutput: test.deleteItemOutput,
	// 			deleteItemError:  test.deleteItemError,
	// 		},
	// 	}
	//
	// 	if err := c.Remove(test.table, test.items...); err != nil && err.Error() != test.err {
	// 		t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
	// 	}
	// }
}
