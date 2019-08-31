package storage

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type s3Client interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	GetObjectRequest(input *s3.GetObjectInput) (req *request.Request, output *s3.GetObjectOutput)
	ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// Storage provides helper methods for persisting/retrieving files
type Storage interface {
	PutFile(int, int, int, int, io.Reader) error
	GetPaths() ([]string, error)
}

// Client implements the S3 interface
type Client struct {
	s3 s3Client
}

// New generates a S3 implementation with an active client
func New() Storage {
	return &Client{
		s3: s3.New(session.New()),
	}
}

// PutFile persists a JSON file in S3
func (c *Client) PutFile(year, month, day, hour int, file io.Reader) error {
	key := fmt.Sprintf("%d/%02d/%02d/%02d/count/%s", year, month, day, hour, uuid.New().String()+"-count.json")

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(file),
		Bucket: aws.String("comana"),
		Key:    aws.String(key),
	}

	_, err := c.s3.PutObject(input)
	if err != nil {
		return fmt.Errorf("error putting file: %s", err.Error())
	}

	return nil
}

var listFiles = func(client s3Client, year, month string, objects *[]*s3.Object) error {
	output, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("comana"),
		Prefix: aws.String(year + "/" + month),
	})
	if err != nil {
		return fmt.Errorf("error listing %s/%s files: %s", year, month, err.Error())
	}

	*objects = append(*objects, output.Contents...)
	return nil
}

var getFile = func(client s3Client, key string) (io.Reader, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String("comana"),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("error getting object %s: %s", key, err.Error())
	}

	return result.Body, nil
}

// GetPaths retrieves paths for files stored in S3
func (c *Client) GetPaths() ([]string, error) {
	current := time.Now()
	year := current.Year()
	month := int(current.Month())

	objects := []*s3.Object{}

	if err := listFiles(c.s3, strconv.Itoa(year), "", &objects); err != nil {
		return nil, fmt.Errorf("error listing files: %s", err.Error())
	}

	for i := 1; i <= 12-month; i++ {
		if err := listFiles(c.s3, strconv.Itoa(year-1), strconv.Itoa(i), &objects); err != nil {
			return nil, fmt.Errorf("error listing files: %s", err.Error())
		}
	}

	paths := []string{}
	for _, object := range objects {
		req, _ := c.s3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String("comana"),
			Key:    aws.String(*object.Key),
		})

		if req == nil {
			return nil, errors.New("error creating get object request")
		}

		signedURL, _ := req.Presign(15 * time.Minute)
		paths = append(paths, signedURL)
	}

	return paths, nil
}
