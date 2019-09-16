package query

import (
	"errors"
	"io"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/aws/aws-lambda-go/events"
)

type mockTable struct {
	getOutput []string
	getCheck  bool
	getErr    error
	addErr    error
	removeErr error
}

func (mock *mockTable) Get(table string, items ...string) ([]string, bool, error) {
	return mock.getOutput, mock.getCheck, mock.getErr
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
	deleteErr   error
}

func (mock *mockStorage) PutFile(key string, file io.Reader) error {
	return mock.putErr
}

func (mock *mockStorage) GetFile(key string) (io.Reader, error) {
	return mock.getOutput, mock.getErr
}

func (mock *mockStorage) GetPaths() ([]string, error) {
	return mock.pathsOutput, mock.pathsErr
}

func (mock *mockStorage) DeleteFile(key string) error {
	return mock.deleteErr
}

func TestCreate(t *testing.T) {
	tests := []struct {
		desc     string
		req      events.APIGatewayProxyRequest
		getCheck bool
		getErr   error
		addErr   error
		putErr   error
		status   int
		err      string
	}{
		{
			desc: "table get error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getCheck: false,
			getErr:   errors.New("mock table get error"),
			addErr:   nil,
			putErr:   nil,
			status:   500,
			err:      "error getting query table: mock table get error",
		},
		{
			desc: "table add error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getCheck: false,
			getErr:   nil,
			addErr:   errors.New("mock table add error"),
			putErr:   nil,
			status:   500,
			err:      "error creating query: mock table add error",
		},
		{
			desc: "table add error",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getCheck: false,
			getErr:   nil,
			addErr:   nil,
			putErr:   errors.New("mock storage put error"),
			status:   500,
			err:      "error putting query file: mock storage put error",
		},
		{
			desc: "table exists",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getCheck: true,
			getErr:   nil,
			addErr:   nil,
			putErr:   nil,
			status:   200,
			err:      "",
		},
		{
			desc: "successful invocation",
			req: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"query_name": "test-name",
				},
			},
			getCheck: false,
			getErr:   nil,
			addErr:   nil,
			putErr:   nil,
			status:   200,
			err:      "",
		},
	}

	for _, test := range tests {
		tbl := &mockTable{
			getCheck: test.getCheck,
			getErr:   test.getErr,
			addErr:   test.addErr,
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

type mockBQ struct {
	bqOutput *bigquery.Query
}

func (mock *mockBQ) Query(query string) *bigquery.Query {
	return mock.bqOutput
}

func TestRun(t *testing.T) {
	tests := []struct {
		desc        string
		pathsOutput []string
		pathsErr    error
		getOutput   io.Reader
		getErr      error
		bqOutput    *bigquery.Query
		subOutput   []string
		subErr      error
		status      int
		err         string
	}{
		{
			desc:        "get paths error",
			pathsOutput: nil,
			pathsErr:    errors.New("mock paths error"),
			getOutput:   nil,
			getErr:      nil,
			bqOutput:    nil,
			subOutput:   nil,
			subErr:      nil,
			status:      500,
			err:         "error listing query files: mock paths error",
		},
		{
			desc:        "get files error",
			pathsOutput: []string{"test-query"},
			pathsErr:    nil,
			getOutput:   nil,
			getErr:      errors.New("mock files error"),
			bqOutput:    nil,
			subOutput:   nil,
			subErr:      nil,
			status:      500,
			err:         "error getting query file: mock files error",
		},
	}

	for _, test := range tests {
		stg := &mockStorage{
			pathsOutput: test.pathsOutput,
			pathsErr:    test.pathsErr,
			getOutput:   test.getOutput,
			getErr:      test.getErr,
		}

		bq := &mockBQ{
			bqOutput: test.bqOutput,
		}

		tbl := &mockTable{
			getOutput: test.subOutput,
			getErr:    test.subErr,
		}

		resp, err := Run(bq, stg, tbl)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		desc      string
		req       events.APIGatewayProxyRequest
		getCheck  bool
		getErr    error
		deleteErr error
		removeErr error
		status    int
		err       string
	}{
		{
			desc: "incorrect secret",
			req: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"COMANA_SECRET": "wrong-test-secret",
				},
			},
			getCheck:  false,
			getErr:    nil,
			deleteErr: nil,
			removeErr: nil,
			status:    500,
			err:       "incorrect secret received: wrong-test-secret",
		},
		{
			desc: "get file error",
			req: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"COMANA_SECRET": "test-secret",
				},
				Body: `{"query": "test-query"}`,
			},
			getCheck:  false,
			getErr:    errors.New("mock get error"),
			deleteErr: nil,
			removeErr: nil,
			status:    500,
			err:       "incorrect getting query: mock get error",
		},
		{
			desc: "delete file error",
			req: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"COMANA_SECRET": "test-secret",
				},
				Body: `{"query": "test-query"}`,
			},
			getCheck:  true,
			getErr:    nil,
			deleteErr: errors.New("mock delete error"),
			removeErr: nil,
			status:    500,
			err:       "incorrect deleting query file: mock delete error",
		},
		{
			desc: "delete item error",
			req: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"COMANA_SECRET": "test-secret",
				},
				Body: `{"query": "test-query"}`,
			},
			getCheck:  true,
			getErr:    nil,
			deleteErr: nil,
			removeErr: errors.New("mock delete error"),
			status:    500,
			err:       "incorrect removing query item: mock delete error",
		},
		{
			desc: "successful invocation",
			req: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"COMANA_SECRET": "test-secret",
				},
				Body: `{"query": "test-query"}`,
			},
			getCheck:  true,
			getErr:    nil,
			deleteErr: nil,
			removeErr: nil,
			status:    200,
			err:       "",
		},
	}

	for _, test := range tests {
		os.Setenv("COMANA_SECRET", "test-secret")

		tbl := &mockTable{
			getOutput: nil,
			getCheck:  test.getCheck,
			getErr:    test.getErr,
			removeErr: test.removeErr,
		}

		stg := &mockStorage{
			deleteErr: test.deleteErr,
		}

		resp, err := Delete(test.req, tbl, stg)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}
