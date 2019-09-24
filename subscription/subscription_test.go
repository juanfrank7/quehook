package subscription

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_createResponse(t *testing.T) {
	tests := []struct {
		desc string
		code int
		msg  string
	}{
		{
			desc: "error message",
			code: 500,
			msg:  "failure",
		},
		{
			desc: "success message",
			code: 200,
			msg:  "success",
		},
	}

	for _, test := range tests {
		resp, _ := createResponse(test.code, test.msg)

		if resp.Body != test.msg {
			t.Errorf("description: %s, body received: %s, expected: %s", test.desc, resp.Body, test.msg)
		}

		if resp.StatusCode != test.code {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.code)
		}
	}
}

type mockTable struct {
	getOutput []string
	getErr    error
	addErr    error
	removeErr error
}

func (mock *mockTable) Get(table string, key, attribute string) ([]string, error) {
	return mock.getOutput, mock.getErr
}

func (mock *mockTable) Add(table string, items ...string) error {
	return mock.addErr
}

func (mock *mockTable) Remove(table string, key, attribute string) error {
	return mock.removeErr
}

func TestSubscribe(t *testing.T) {
	tests := []struct {
		desc      string
		req       events.APIGatewayProxyRequest
		getOutput []string
		getErr    error
		addErr    error
		code      int
		err       string
	}{
		{
			desc: "unmarshalling request body error",
			req: events.APIGatewayProxyRequest{
				Body: `["c-3po", "r2-d2"]`,
			},
			getOutput: nil,
			getErr:    nil,
			addErr:    nil,
			code:      400,
			err:       "error parsing request: json: cannot unmarshal array into Go value of type subscription.sub",
		},
		{
			desc: "incorrect url target",
			req: events.APIGatewayProxyRequest{
				Body: `{"subscriber_target":"liesofthejedi.com"}`,
			},
			getOutput: nil,
			getErr:    nil,
			addErr:    nil,
			code:      400,
			err:       "error parsing url: parse liesofthejedi.com: invalid URI for request",
		},
		{
			desc: "error getting query",
			req: events.APIGatewayProxyRequest{
				Body: `{"subscriber_target":"https://tragedyofdarthplagueisthewise.com"}`,
			},
			getOutput: nil,
			getErr:    errors.New("mock get error"),
			addErr:    nil,
			code:      500,
			err:       "error getting query: mock get error",
		},
		{
			desc: "error getting query output",
			req: events.APIGatewayProxyRequest{
				Body: `{"subscriber_target":"https://jediarchives.com/kamino"}`,
			},
			getOutput: []string{},
			getErr:    nil,
			addErr:    nil,
			code:      500,
			err:       "no output returned",
		},
		{
			desc: "add subscriber error",
			req: events.APIGatewayProxyRequest{
				Body: `{"query_name":"kamino","subscriber_target":"https://kenobi.io","subscriber_name":"obi-wan","subscriber_email":"kenobi@jedi.ord"}`,
			},
			getOutput: []string{"nu"},
			getErr:    nil,
			addErr:    errors.New("mock add error"),
			code:      500,
			err:       "error adding subscriber: mock add error",
		},
		{
			desc: "successful invocation",
			req: events.APIGatewayProxyRequest{
				Body: `{"query_name":"kamino","subscriber_target":"https://kenobi.io","subscriber_name":"obi-wan","subscriber_email":"kenobi@jedi.ord"}`,
			},
			getOutput: []string{"dexter-jetster"},
			getErr:    nil,
			addErr:    nil,
			code:      200,
			err:       "",
		},
	}

	for _, test := range tests {
		tbl := &mockTable{
			getOutput: test.getOutput,
			getErr:    test.getErr,
			addErr:    test.addErr,
		}

		resp, err := Subscribe(test.req, tbl)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.code {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.code)
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	tests := []struct {
		desc      string
		req       events.APIGatewayProxyRequest
		getOutput []string
		getErr    error
		removeErr error
		code      int
		err       string
	}{
		{
			desc: "error parsing request",
			req: events.APIGatewayProxyRequest{
				Body: `["incorrect-input"]`,
			},

			getOutput: nil,
			getErr:    errors.New("mock get error"),
			removeErr: nil,
			code:      400,
			err:       "error parsing request: json: cannot unmarshal array into Go value of type subscription.sub",
		},
		{
			desc: "error getting subscriber",
			req: events.APIGatewayProxyRequest{
				Body: `{"query_name":"droids"}`,
			},
			getOutput: []string{""},
			getErr:    errors.New("mock get error"),
			removeErr: nil,
			code:      500,
			err:       "error getting query: mock get error",
		},
		{
			desc: "error removing subscriber",
			req: events.APIGatewayProxyRequest{
				Body: `{"query_name":"droids"}`,
			},
			getOutput: []string{"not-them"},
			getErr:    nil,
			removeErr: errors.New("mock remove error"),
			code:      500,
			err:       "error removing subscriber: mock remove error",
		},
		{
			desc: "successful invocation",
			req: events.APIGatewayProxyRequest{
				Body: `{"query_name":"droids"}`,
			},
			getOutput: []string{"move along"},
			getErr:    nil,
			removeErr: nil,
			code:      200,
			err:       "",
		},
	}

	for _, test := range tests {
		tbl := &mockTable{
			getOutput: test.getOutput,
			getErr:    test.getErr,
			removeErr: test.removeErr,
		}

		resp, err := Unsubscribe(test.req, tbl)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.code {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.code)
		}
	}
}
