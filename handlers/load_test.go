package handlers

import (
	"errors"
	"testing"
)

func TestLoadData(t *testing.T) {
	tests := []struct {
		desc        string
		getPathsOut []string
		getPathsErr error
		status      int
		err         string
	}{
		{
			desc:        "get paths error",
			getPathsOut: nil,
			getPathsErr: errors.New("get paths error"),
			status:      500,
			err:         "get paths error",
		},
		{
			desc:        "successful invocation",
			getPathsOut: []string{},
			getPathsErr: nil,
			status:      200,
			err:         "",
		},
	}

	for _, test := range tests {
		s := &mockStorage{
			getPathsOut: test.getPathsOut,
			getPathsErr: test.getPathsErr,
		}

		resp, err := LoadData(s)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}
