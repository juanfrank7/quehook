package storage

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Client interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

// Storage provides helper methods for persisting/retrieving files
type Storage interface {
	PutFile(string, io.Reader) error
	GetFile(string) (io.Reader, error)
	GetPaths() ([]string, error)
	DeleteFile(string) error
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
func (c *Client) PutFile(key string, file io.Reader) error {
	_, err := c.s3.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(file),
		Bucket: aws.String("comana"),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("error putting file: %s", err.Error())
	}

	return nil
}

// GetFile retrieves a given file stored in S3
func (c *Client) GetFile(key string) (io.Reader, error) {
	result, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("comana"),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting object %s: %s", key, err.Error())
	}

	return result.Body, nil
}

// GetPaths retrieves paths for files stored in S3
func (c *Client) GetPaths() ([]string, error) {
	output, err := c.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("comana"),
		Prefix: aws.String("queries"),
	})

	if err != nil {
		return nil, fmt.Errorf("error listing files: %s", err.Error())
	}

	paths := []string{}
	for _, object := range output.Contents {
		paths = append(paths, *object.Key)
	}

	return paths, nil
}

// DeleteFile removes a query file from S3
func (c *Client) DeleteFile(key string) error {
	_, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("comana"),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("error deleting file: %s", err.Error())
	}

	return nil
}
