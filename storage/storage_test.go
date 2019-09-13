package storage

import (
	"errors"
	"io"
	// 	"io/ioutil"
	// 	"net/http"
	// 	"net/url"
	"strings"
	"testing"
	// 	"time"

	// 	"github.com/aws/aws-sdk-go/aws"
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
	listObjectsOutput  *s3.ListObjectsV2Output
	listObjectsErr     error
	putObjectOutput    *s3.PutObjectOutput
	putObjectErr       error
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

		if err := c.PutFile("episode i", test.file); err != nil && err.Error() != test.err {
			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
		}
	}
}

// func TestGetPaths(t *testing.T) {
// 	tests := []struct {
// 		desc               string
// 		listErr            error
// 		getObjectReq       *request.Request
// 		getObjectReqOutput *s3.GetObjectOutput
// 		outputLen          int
// 		err                string
// 	}{
// 		{
// 			desc:               "list files error",
// 			listErr:            errors.New("listing error"),
// 			getObjectReq:       nil,
// 			getObjectReqOutput: nil,
// 			outputLen:          0,
// 			err:                "error listing files: listing error",
// 		},
// 		{
// 			desc:               "get object request error",
// 			listErr:            nil,
// 			getObjectReq:       nil,
// 			getObjectReqOutput: nil,
// 			outputLen:          0,
// 			err:                "error creating get object request",
// 		},
// 		{
// 			desc:               "request creation error",
// 			listErr:            nil,
// 			getObjectReq:       nil,
// 			getObjectReqOutput: nil,
// 			outputLen:          0,
// 			err:                "error creating get object request",
// 		},
// 		{
// 			desc:    "successful invocation",
// 			listErr: nil,
// 			getObjectReq: &request.Request{
// 				Operation: &request.Operation{},
// 				HTTPRequest: &http.Request{
// 					URL: &url.URL{},
// 				},
// 			},
// 			getObjectReqOutput: nil,
// 			outputLen:          12 - int(time.Now().Month()) + 1, // Due to time.Now() used in GetPaths method
// 			err:                "",
// 		},
// 	}
//
// 	for _, test := range tests {
// 		c := &Client{
// 			s3: &storageMock{
// 				getObjectReq:       test.getObjectReq,
// 				getObjectReqOutput: test.getObjectReqOutput,
// 			},
// 		}
//
// 		listFiles = func(client s3Client, year, month string, objects *[]*s3.Object) error {
// 			*objects = append(*objects, &s3.Object{
// 				Key: aws.String("test-key"),
// 			})
// 			return test.listErr
// 		}
//
// 		output, err := c.GetPaths()
// 		if err != nil && err.Error() != test.err {
// 			t.Errorf("description: %s, error received: %s, expected: %s", test.desc, err.Error(), test.err)
// 		}
//
// 		if output != nil && len(output) != test.outputLen {
// 			t.Errorf("description: %s, length received: %d, expected: %d", test.desc, len(output), test.outputLen)
// 		}
// 	}
// }
