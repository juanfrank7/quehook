package handlers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
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
			"",
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
			desc: "successful invocation with single value",
			scnr: bufio.NewScanner(
				strings.NewReader(`{"type": "test-event", "repo":{"name": "test-repo"}}`),
			),
			rdr: strings.NewReader("test-reader"),
			err: "",
		},
		{
			desc: "successful invocation with multiple values",
			scnr: bufio.NewScanner(
				strings.NewReader(`{"type": "test-event", "repo":{"name": "test-repo"}}\n{"type": "test-event", "repo":{"name": "test-repo"}}`),
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
			err:    "source must be cloudwatch event or backfill",
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
