package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
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
	getObjectReq       *request.Request
	getObjectReqOutput *s3.GetObjectOutput
	listObjectsOutput  *s3.ListObjectsOutput
	listObjectsErr     error
	putObjectOutput    *s3.PutObjectOutput
	putObjectErr       error
}

func (mock *storageMock) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return mock.getObjectOutput, mock.getObjectError
}

func (mock *storageMock) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return mock.listObjectsOutput, mock.listObjectsErr
}

func (mock *storageMock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return mock.putObjectOutput, mock.putObjectErr
}

func (mock *storageMock) GetObjectRequest(input *s3.GetObjectInput) (req *request.Request, output *s3.GetObjectOutput) {
	return mock.getObjectReq, mock.getObjectReqOutput
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

		if err := c.PutFile(test.file); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func Test_listFiles(t *testing.T) {
	tests := []struct {
		desc       string
		listOutput *s3.ListObjectsOutput
		listErr    error
		length     int
		err        string
	}{
		{
			desc:       "s3 client error",
			listOutput: nil,
			listErr:    errors.New("mock storage error"),
			length:     0,
			err:        "error listing 1977/5 files: mock storage error",
		},
		{
			desc: "successful invocation",
			listOutput: &s3.ListObjectsOutput{
				Contents: []*s3.Object{
					&s3.Object{
						Key: aws.String("a-new-hope"),
					},
				},
			},
			listErr: nil,
			length:  1,
			err:     "",
		},
	}

	for _, test := range tests {
		c := &storageMock{
			listObjectsOutput: test.listOutput,
			listObjectsErr:    test.listErr,
		}

		objects := &[]*s3.Object{}
		if err := listFiles(c, "1977", "5", objects); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if len(*objects) != test.length {
			t.Errorf("description: %s, output length received: %d, expected: %d", test.desc, len(*objects), test.length)
		}
	}
}

func Test_getFile(t *testing.T) {
	tests := []struct {
		desc      string
		getOutput *s3.GetObjectOutput
		getErr    error
		err       string
	}{
		{
			desc:      "s3 client error",
			getOutput: nil,
			getErr:    errors.New("mock storage error"),
			err:       "error getting object key: mock storage error",
		},
		{
			desc: "successful invocation",
			getOutput: &s3.GetObjectOutput{
				Body: ioutil.NopCloser(strings.NewReader("test")),
			},
			getErr: nil,
			err:    "",
		},
	}

	for _, test := range tests {
		c := &storageMock{
			getObjectOutput: test.getOutput,
			getObjectError:  test.getErr,
		}

		if _, err := getFile(c, "key"); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

func TestGetFiles(t *testing.T) {
	tests := []struct {
		desc      string
		listErr   error
		getOutput io.Reader
		getErr    error
		outputLen int
		err       string
	}{
		{
			desc:      "list files error",
			listErr:   errors.New("listing error"),
			getOutput: strings.NewReader("mock"),
			getErr:    nil,
			outputLen: 0,
			err:       "error listing files: listing error",
		},
		{
			desc:      "get file error",
			listErr:   nil,
			getOutput: strings.NewReader("mock"),
			getErr:    errors.New("get error"),
			outputLen: 0,
			err:       "error getting object test-key: get error",
		},
		{
			desc:      "successful invocation",
			listErr:   nil,
			getOutput: strings.NewReader("mock"),
			getErr:    nil,
			outputLen: 1, // there is only one key in the map
			err:       "",
		},
	}

	for _, test := range tests {
		c := &Client{}

		listFiles = func(client s3Client, year, month string, objects *[]*s3.Object) error {
			*objects = append(*objects, &s3.Object{
				Key: aws.String("test-key"),
			})
			return test.listErr
		}

		getFile = func(client s3Client, key string) (io.Reader, error) {
			return test.getOutput, test.getErr
		}

		output, err := c.GetFiles()
		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if output != nil && len(output) != test.outputLen {
			t.Errorf("description: %s, length received: %d, expected: %d", test.desc, len(output), test.outputLen)
		}
	}
}

func TestGetPaths(t *testing.T) {
	tests := []struct {
		desc               string
		listErr            error
		getObjectReq       *request.Request
		getObjectReqOutput *s3.GetObjectOutput
		outputLen          int
		err                string
	}{
		{
			desc:               "list files error",
			listErr:            errors.New("listing error"),
			getObjectReq:       nil,
			getObjectReqOutput: nil,
			outputLen:          0,
			err:                "error listing files: listing error",
		},
		{
			desc:               "get object request error",
			listErr:            nil,
			getObjectReq:       nil,
			getObjectReqOutput: nil,
			outputLen:          0,
			err:                "error creating get object request",
		},
		{
			desc:               "request creation error",
			listErr:            nil,
			getObjectReq:       nil,
			getObjectReqOutput: nil,
			outputLen:          0,
			err:                "error creating get object request",
		},
		{
			desc:    "successful invocation",
			listErr: nil,
			getObjectReq: &request.Request{
				Operation: &request.Operation{},
				HTTPRequest: &http.Request{
					URL: &url.URL{},
				},
			},
			getObjectReqOutput: nil,
			outputLen:          5,
			err:                "",
		},
	}

	for _, test := range tests {
		c := &Client{
			s3: &storageMock{
				getObjectReq:       test.getObjectReq,
				getObjectReqOutput: test.getObjectReqOutput,
			},
		}

		listFiles = func(client s3Client, year, month string, objects *[]*s3.Object) error {
			*objects = append(*objects, &s3.Object{
				Key: aws.String("test-key"),
			})
			return test.listErr
		}

		output, err := c.GetPaths()
		if err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}

		if output != nil && len(output) != test.outputLen {
			t.Errorf("description: %s, length received: %d, expected: %d", test.desc, len(output), test.outputLen)
		}
	}
}
