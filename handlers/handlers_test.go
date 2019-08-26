package handlers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_download(t *testing.T) {
	tests := []struct {
		desc string
		url  string
		err  string
	}{
		{
			"bad url",
			"https://temp.com",
			"Get https://temp.com: dial tcp 23.23.86.44:443: connect: connection refused",
		},
		{
			"good url",
			"https://data.gharchive.org/2019-01-01-15.json.gz",
			"asdfasdfas",
		},
	}

	for _, test := range tests {
		file, err := download(test.url)
		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if err == nil && file == nil {
			t.Errorf("description: %s, no output file found", test.desc)
		}
	}
}

func Test_unzip(t *testing.T) {
	tests := []struct {
		desc  string
		input []byte
		err   string
	}{
		{
			desc:  "successful invocation",
			input: []byte(`{"test-key": "test-value"}`),
			err:   "",
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer

		w := gzip.NewWriter(&buf)
		w.Header = gzip.Header{
			Name:    "test.json.gz",
			Comment: "",
			Extra:   []byte{},
			ModTime: time.Now(),
			OS:      byte(255),
		}
		_, err := w.Write(test.input)
		if err != nil {
			log.Fatalf("error creating test gzip resource %s", err.Error())
		}

		_, err = unzip(buf.Bytes())

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		desc string
		scnr *bufio.Scanner
		rdr  io.Reader
		err  string
	}{
		{
			desc: "successful invocation",
			scnr: bufio.NewScanner(
				strings.NewReader(`{"type": "test-event", "repo":{"name": "test-repo"}}`),
			),
			rdr: strings.NewReader("test-reader"),
			err: "",
		},
	}

	for _, test := range tests {
		_, err := parse(test.scnr)
		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

type mockStorage struct {
	putFileErr  error
	getFilesOut map[string]io.Reader
	getFilesErr error
	getPathsOut []string
	getPathsErr error
}

func (m *mockStorage) PutFile(io.Reader) error {
	return m.putFileErr
}

func (m *mockStorage) GetFiles() (map[string]io.Reader, error) {
	return m.getFilesOut, m.getFilesErr
}

func (m *mockStorage) GetPaths() ([]string, error) {
	return m.getPathsOut, m.getPathsErr
}

func TestSaveData(t *testing.T) {
	tests := []struct {
		desc   string
		src    string
		dwn    func(url string) ([]byte, error)
		uzp    func([]byte) (*bufio.Scanner, error)
		prs    func(s *bufio.Scanner) (io.Reader, error)
		dbErr  error
		status int
		err    string
	}{
		{
			desc: "incorrect source",
			src:  "not-source",
			dwn: func(string) ([]byte, error) {
				return nil, errors.New("download error")
			},
			uzp:    nil,
			prs:    nil,
			dbErr:  nil,
			status: 500,
			err:    "request source must be cloudwatch event",
		},
		{
			desc: "archive download error",
			src:  "aws.events",
			dwn: func(string) ([]byte, error) {
				return nil, errors.New("download error")
			},
			uzp:    nil,
			prs:    nil,
			dbErr:  nil,
			status: 500,
			err:    "download error",
		},
		{
			desc: "archive unzip error",
			src:  "aws.events",
			dwn: func(string) ([]byte, error) {
				return nil, nil
			},
			uzp: func([]byte) (*bufio.Scanner, error) {
				return nil, errors.New("unzip error")
			},
			prs:    nil,
			dbErr:  nil,
			status: 500,
			err:    "unzip error",
		},
		{
			desc: "archive parse error",
			src:  "aws.events",
			dwn: func(string) ([]byte, error) {
				return nil, nil
			},
			uzp: func([]byte) (*bufio.Scanner, error) {
				return nil, nil
			},
			prs: func(s *bufio.Scanner) (io.Reader, error) {
				return nil, errors.New("parse error")
			},
			dbErr:  nil,
			status: 500,
			err:    "parse error",
		},
		{
			desc: "put file error",
			src:  "aws.events",
			dwn: func(string) ([]byte, error) {
				return nil, nil
			},
			uzp: func([]byte) (*bufio.Scanner, error) {
				return nil, nil
			},
			prs: func(s *bufio.Scanner) (io.Reader, error) {
				return strings.NewReader("test"), nil
			},
			dbErr:  errors.New("put file error"),
			status: 500,
			err:    "put file error",
		},
		{
			desc: "successful invocation",
			src:  "aws.events",
			dwn: func(string) ([]byte, error) {
				return nil, nil
			},
			uzp: func([]byte) (*bufio.Scanner, error) {
				return nil, nil
			},
			prs: func(s *bufio.Scanner) (io.Reader, error) {
				return strings.NewReader("test"), nil
			},
			dbErr:  nil,
			status: 200,
			err:    "",
		},
	}

	for _, test := range tests {
		s := &mockStorage{
			putFileErr: test.dbErr,
		}

		download = test.dwn
		unzip = test.uzp
		parse = test.prs

		req := Request{
			Source: test.src,
		}

		resp, err := SaveData(req, s)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}

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

func TestBackfillData(t *testing.T) {
	tests := []struct {
		desc        string
		secret      string
		body        string
		downloadErr error
		unzipOutput *bufio.Scanner
		unzipErr    error
		parseOutput io.Reader
		parseErr    error
		putFileErr  error
		status      int
		err         string
	}{
		{
			desc:        "incorrect request secret",
			secret:      "test-secret-failure",
			body:        "",
			downloadErr: nil,
			unzipOutput: nil,
			unzipErr:    nil,
			parseOutput: strings.NewReader("test"),
			parseErr:    nil,
			putFileErr:  nil,
			status:      500,
			err:         "incorrect secret received: test-secret-failure",
		},
		{
			desc:        "download function error",
			secret:      "test-secret",
			body:        `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			downloadErr: errors.New("download error"),
			unzipOutput: nil,
			unzipErr:    nil,
			parseOutput: strings.NewReader("test"),
			parseErr:    nil,
			putFileErr:  nil,
			status:      500,
			err:         "download error",
		},
		{
			desc:        "unzip function error",
			secret:      "test-secret",
			body:        `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			downloadErr: nil,
			unzipOutput: nil,
			unzipErr:    errors.New("unzip error"),
			parseOutput: strings.NewReader("test"),
			parseErr:    nil,
			putFileErr:  nil,
			status:      500,
			err:         "unzip error",
		},
		{
			desc:        "parse function error",
			secret:      "test-secret",
			body:        `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			downloadErr: nil,
			unzipOutput: nil,
			unzipErr:    nil,
			parseOutput: strings.NewReader("test"),
			parseErr:    errors.New("parse error"),
			putFileErr:  nil,
			status:      500,
			err:         "parse error",
		},
		{
			desc:        "parse function error",
			secret:      "test-secret",
			body:        `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			downloadErr: nil,
			unzipOutput: nil,
			unzipErr:    nil,
			parseOutput: strings.NewReader("test"),
			parseErr:    nil,
			putFileErr:  errors.New("put file error"),
			status:      500,
			err:         "put file error",
		},
		{
			desc:        "parse function error",
			secret:      "test-secret",
			body:        `{"year": 1977, "month": 5, "startDay": 25, "endDay": 25}`,
			downloadErr: nil,
			unzipOutput: nil,
			unzipErr:    nil,
			parseOutput: strings.NewReader("test"),
			parseErr:    nil,
			putFileErr:  nil,
			status:      200,
			err:         "",
		},
	}

	for _, test := range tests {
		os.Setenv("COMANA_SECRET", "test-secret")

		s := &mockStorage{
			putFileErr: test.putFileErr,
		}

		download = func(url string) ([]byte, error) {
			return nil, test.downloadErr
		}

		unzip = func([]byte) (*bufio.Scanner, error) {
			return test.unzipOutput, test.unzipErr
		}

		parse = func(s *bufio.Scanner) (io.Reader, error) {
			return test.parseOutput, test.parseErr
		}

		r := Request{
			Headers: map[string]string{
				"COMANA_SECRET": test.secret,
			},
			Body: test.body,
		}

		resp, err := BackfillData(r, s)

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if resp.StatusCode != test.status {
			t.Errorf("description: %s, status received: %d, expected: %d", test.desc, resp.StatusCode, test.status)
		}
	}
}
