package table

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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
	tests := []struct {
		desc          string
		table         string
		key           string
		getItemOutput *dynamodb.BatchGetItemOutput
		getItemError  error
		output        []string
		err           string
	}{
		{
			desc:          "get item error",
			table:         "queries",
			key:           "query",
			getItemOutput: nil,
			getItemError:  errors.New("mock get error"),
			output:        nil,
			err:           "get item error: mock get error",
		},
		{
			desc:  "successful subscribers invocation",
			table: "subscribers",
			key:   "query",
			getItemOutput: &dynamodb.BatchGetItemOutput{
				Responses: map[string][]map[string]*dynamodb.AttributeValue{
					"subscribers": []map[string]*dynamodb.AttributeValue{
						map[string]*dynamodb.AttributeValue{
							"query_name": {
								S: aws.String("test-query"),
							},
							"subscriber_target": {
								S: aws.String("test-target"),
							},
						},
					},
				},
			},
			getItemError: nil,
			output: []string{
				"test-query",
			},
			err: "",
		},
		{
			desc:  "successful queries invocation",
			table: "queries",
			key:   "key",
			getItemOutput: &dynamodb.BatchGetItemOutput{
				Responses: map[string][]map[string]*dynamodb.AttributeValue{
					"queries": []map[string]*dynamodb.AttributeValue{
						map[string]*dynamodb.AttributeValue{
							"query_name": {
								S: aws.String("test-query"),
							},
						},
					},
				},
			},
			getItemError: nil,
			output: []string{
				"test-query",
			},
			err: "",
		},
	}

	for _, test := range tests {
		c := &Client{
			dynamodb: &tableMock{
				getItemOutput: test.getItemOutput,
				getItemError:  test.getItemError,
			},
		}

		output, err := c.Get(test.table, test.key)

		if len(output) != len(test.output) {
			t.Errorf("description: %s, output received: %d, expected: %d", test.desc, len(output), len(test.output))
		}

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		desc             string
		table            string
		key              string
		getItemOutput    *dynamodb.BatchGetItemOutput
		getItemError     error
		deleteItemOutput *dynamodb.DeleteItemOutput // kept for possible method expansion
		deleteItemError  error
		err              string
	}{
		{
			desc:             "delete queries error",
			table:            "queries",
			key:              "query",
			getItemOutput:    nil,
			getItemError:     nil,
			deleteItemOutput: nil,
			deleteItemError:  errors.New("mock delete error"),
			err:              "delete item error: mock delete error",
		},
		{
			desc:             "delete queries successful invocation",
			table:            "queries",
			key:              "query",
			getItemOutput:    nil,
			getItemError:     nil,
			deleteItemOutput: nil,
			deleteItemError:  nil,
			err:              "",
		},
		{
			desc:             "delete subscribers get batch error",
			table:            "subscribers",
			key:              "query",
			getItemOutput:    nil,
			getItemError:     errors.New("mock get error"),
			deleteItemOutput: nil,
			deleteItemError:  nil,
			err:              "get item error: mock get error",
		},
		{
			desc:  "delete subscribers delete error",
			table: "subscribers",
			key:   "query",
			getItemOutput: &dynamodb.BatchGetItemOutput{
				Responses: map[string][]map[string]*dynamodb.AttributeValue{
					"subscribers": []map[string]*dynamodb.AttributeValue{
						map[string]*dynamodb.AttributeValue{
							"subscriber_email": &dynamodb.AttributeValue{
								S: aws.String("anakin@skywalker.com"),
							},
						},
					},
				},
			},
			getItemError:     nil,
			deleteItemOutput: nil,
			deleteItemError:  errors.New("mock delete error"),
			err:              "delete item error: mock delete error",
		},
		{
			desc:  "delete subscribers successful invocation",
			table: "subscribers",
			key:   "query",
			getItemOutput: &dynamodb.BatchGetItemOutput{
				Responses: map[string][]map[string]*dynamodb.AttributeValue{
					"subscribers": []map[string]*dynamodb.AttributeValue{
						map[string]*dynamodb.AttributeValue{
							"subscriber_email": &dynamodb.AttributeValue{
								S: aws.String("anakin@skywalker.com"),
							},
						},
					},
				},
			},
			getItemError:     nil,
			deleteItemOutput: nil,
			deleteItemError:  nil,
			err:              "",
		},
	}

	for _, test := range tests {
		c := &Client{
			dynamodb: &tableMock{
				getItemOutput:    test.getItemOutput,
				getItemError:     test.getItemError,
				deleteItemOutput: test.deleteItemOutput,
				deleteItemError:  test.deleteItemError,
			},
		}

		if err := c.Remove(test.table, test.key); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}
