package storage

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Error("description: error creating new storage implementation")
	}
}

type storageMock struct {
	getObjectOutput    *s3.GetObjectOutput
	getObjectError     error
	listObjectsOutput  *s3.ListObjectsV2Output
	listObjectsErr     error
	putObjectOutput    *s3.PutObjectOutput
	putObjectErr       error
	deleteObjectOutput *s3.DeleteObjectOutput
	deleteObjectErr    error
}

func (mock *storageMock) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return mock.getObjectOutput, mock.getObjectError
}

func (mock *storageMock) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return mock.listObjectsOutput, mock.listObjectsErr
}

func (mock *storageMock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return mock.putObjectOutput, mock.putObjectErr
}

func (mock *storageMock) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return mock.deleteObjectOutput, mock.deleteObjectErr
}

func TestPutFile(t *testing.T) {
	tests := []struct {
		desc          string
		file          io.Reader
		storageOutput *s3.PutObjectOutput
		storageErr    error
		err           string
	}{
		{
			desc:          "s3 client error",
			file:          strings.NewReader("test"),
			storageOutput: &s3.PutObjectOutput{},
			storageErr:    errors.New("mock storage error"),
			err:           "error putting file: mock storage error",
		},
		{
			desc:          "successful invocation",
			file:          strings.NewReader("test"),
			storageOutput: &s3.PutObjectOutput{},
			storageErr:    nil,
			err:           "",
		},
	}

	for _, test := range tests {
		c := &Client{
			s3: &storageMock{
				putObjectOutput: test.storageOutput,
				putObjectErr:    test.storageErr,
			},
		}

		if err := c.PutFile("episode i", test.file); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func TestGetFile(t *testing.T) {
	tests := []struct {
		desc            string
		getObjectOutput *s3.GetObjectOutput
		getObjectError  error
		output          string
		err             string
	}{
		{
			desc:            "get object error",
			getObjectOutput: nil,
			getObjectError:  errors.New("mock get file error"),
			output:          "",
			err:             "error getting object test-key: mock get file error",
		},
		{
			desc: "successful invocation",
			getObjectOutput: &s3.GetObjectOutput{
				Body: ioutil.NopCloser(strings.NewReader("example query")),
			},
			getObjectError: nil,
			output:         "example query",
			err:            "",
		},
	}

	for _, test := range tests {
		c := &Client{
			s3: &storageMock{
				getObjectOutput: test.getObjectOutput,
				getObjectError:  test.getObjectError,
			},
		}

		output, err := c.GetFile("test-key")

		if output != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(output)
			file := buf.String()

			if file != test.output {
				t.Errorf("description: %s, output received: %s, expected: %s", test.desc, file, test.output)
			}
		}

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func TestGetPaths(t *testing.T) {
	tests := []struct {
		desc              string
		listObjectsOutput *s3.ListObjectsV2Output
		listObjectsErr    error
		output            []string
		err               string
	}{
		{
			desc:              "get paths error",
			listObjectsOutput: nil,
			listObjectsErr:    errors.New("mock list error"),
			output:            nil,
			err:               "error listing files: mock list error",
		},
		{
			desc: "successful invocation",
			listObjectsOutput: &s3.ListObjectsV2Output{
				Contents: []*s3.Object{
					&s3.Object{
						Key: aws.String("test-key"),
					},
				},
			},
			output: []string{"test-key"},
			err:    "",
		},
	}

	for _, test := range tests {
		c := &Client{
			s3: &storageMock{
				listObjectsOutput: test.listObjectsOutput,
				listObjectsErr:    test.listObjectsErr,
			},
		}

		output, err := c.GetPaths()

		if output != nil && len(output) != len(test.output) {
			t.Errorf("description: %s, output received: %s, expected: %s", test.desc, output, test.output)
		}

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		desc               string
		deleteObjectOutput *s3.DeleteObjectOutput
		deleteObjectErr    error
		err                string
	}{
		{
			desc:               "delete file error",
			deleteObjectOutput: nil,
			deleteObjectErr:    errors.New("mock delete error"),
			err:                "error deleting file: mock delete error",
		},
		{
			desc:               "successful invocation",
			deleteObjectOutput: nil,
			deleteObjectErr:    nil,
			err:                "",
		},
	}

	for _, test := range tests {
		c := &Client{
			s3: &storageMock{
				deleteObjectOutput: test.deleteObjectOutput,
				deleteObjectErr:    test.deleteObjectErr,
			},
		}

		err := c.DeleteFile("test-key")

		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}
