package query

import (
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type mockTable struct {
	getOutput bool
	getErr    error
	addErr    error
	removeErr error
}

func (mock *mockTable) Get(table string, items ...string) (bool, error) {
	return mock.getOutput, mock.getErr
}

func (mock *mockTable) Add(table string, items ...string) error {
	return mock.addErr
}

func (mock *mockTable) Remove(table string, item ...string) error {
	return mock.removeErr
}

type mockStorage struct {
	putErr      error
	getOutput   io.Reader
	getErr      error
	pathsOutput []string
	pathsErr    error
}

func (mock *mockStorage) PutFile(key string, file io.Reader) error {
	return mock.putErr
}

func (mock *mockStorage) GetFile(string) (io.Reader, error) {
	return mock.getOutput, mock.getErr
}

func (mock *mockStorage) GetPaths() ([]string, error) {
	return mock.pathsOutput, mock.pathsErr
}

func TestCreate(t *testing.T) {
	tests := []struct {
		desc      string
		req       events.APIGatewayProxyRequest
		getOutput bool
		getErr    error
		addErr    error
		putErr    error
		status    int
		err       string
	}{
		{
			desc: "table get error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getOutput: false,
			getErr:    errors.New("mock table get error"),
			addErr:    nil,
			putErr:    nil,
			status:    500,
			err:       "error getting query table: mock table get error",
		},
		{
			desc: "table add error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getOutput: false,
			getErr:    nil,
			addErr:    errors.New("mock table add error"),
			putErr:    nil,
			status:    500,
			err:       "error creating query: mock table add error",
		},
		{
			desc: "table add error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getOutput: false,
			getErr:    nil,
			addErr:    nil,
			putErr:    errors.New("mock storage put error"),
			status:    500,
			err:       "error putting query file: mock storage put error",
		},
		{
			desc: "table exists",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getOutput: true,
			getErr:    nil,
			addErr:    nil,
			putErr:    nil,
			status:    200,
			err:       "",
		},
		{
			desc: "successful invocation",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getOutput: false,
			getErr:    nil,
			addErr:    nil,
			putErr:    nil,
			status:    200,
			err:       "",
		},
	}

	for _, test := range tests {
		tbl := &mockTable{
			getOutput: test.getOutput,
			getErr:    test.getErr,
			addErr:    test.addErr,
		}

		stg := &mockStorage{
			putErr: test.putErr,
		}

		resp, err := Create(test.req, tbl, stg)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}
