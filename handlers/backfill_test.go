package handlers

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestNewInvoke(t *testing.T) {
	i := NewInvoke()
	if i == nil {
		t.Error("description: error creating new invoke implementation")
	}
}

type mockInvoke struct {
	invokeStatus int64
	invokeResp   string
	invokeErr    error
}

func (m *mockInvoke) Invoke(payload []byte) (int64, string, error) {
	return m.invokeStatus, m.invokeResp, m.invokeErr
}

func TestBackfillData(t *testing.T) {
	tests := []struct {
		desc         string
		secret       string
		body         string
		invokeStatus int64
		invokeResp   string
		invokeErr    error
		status       int
		err          string
	}{
		{
			desc:         "incorrect request secret",
			secret:       "test-secret-failure",
			body:         "",
			invokeStatus: 0,
			invokeResp:   "",
			invokeErr:    nil,
			status:       500,
			err:          "incorrect secret received: test-secret-failure",
		},
		{
			desc:         "invoke method error",
			secret:       "test-secret",
			body:         `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			invokeStatus: 500,
			invokeResp:   "invoke-error",
			invokeErr:    errors.New("invoke-error"),
			status:       500,
			err:          "lambda invocation error for https://data.gharchive.org/1977-05-00",
		},
		{
			desc:         "successful invocation",
			secret:       "test-secret",
			body:         `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			invokeStatus: 200,
			invokeResp:   "invoke-success",
			invokeErr:    nil,
			status:       200,
			err:          "",
		},
	}

	for _, test := range tests {
		os.Setenv("COMANA_SECRET", "test-secret")

		i := &mockInvoke{
			invokeStatus: test.invokeStatus,
			invokeResp:   test.invokeResp,
			invokeErr:    test.invokeErr,
		}

		r := Request{
			Headers: map[string]string{
				"COMANA_SECRET": test.secret,
			},
			Body: test.body,
		}

		resp, err := BackfillData(r, i)

		if err != nil && !strings.Contains(err.Error(), test.err) {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}
