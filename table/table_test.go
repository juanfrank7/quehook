package table

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type tableMock struct {
	putItemOutput    *dynamodb.PutItemOutput
	putItemError     error
	getItemOutput    *dynamodb.GetItemOutput
	getItemError     error
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemError  error
}

func (mock *tableMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return mock.putItemOutput, mock.putItemError
}

func (mock *tableMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
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
				"query",
			},
			putItemOutput: nil,
			putItemError:  errors.New("mock put error"),
			err:           "put item error: mock put error",
		},
		{
			desc:  "successful subscribers invocation",
			table: "subscribers",
			items: []string{
				"query",
				"subname",
				"target",
			},
			putItemOutput: nil,
			putItemError:  nil,
			err:           "",
		},
		{
			desc:  "successful queries invocation",
			table: "queries",
			items: []string{
				"query",
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

}

func TestRemove(t *testing.T) {

}
