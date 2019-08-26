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
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// Storage provides helper methods for persisting/retrieving files
type Storage interface {
	PutFile(io.Reader) error
	GetFiles() (map[string]io.Reader, error)
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
func (c *Client) PutFile(file io.Reader) error {
	current := time.Now().Add(time.Hour * -5) // 1 hour prior + 4 hour UTC-EST difference
	year := current.Year()
	month := int(current.Month())
	day := current.Day()
	hour := current.Hour()

	key := fmt.Sprintf("%d/%02d/%02d/%02d/%s", year, month, day, hour, uuid.New().String()+"-count.json")

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
	output, err := client.ListObjects(&s3.ListObjectsInput{
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

// GetFiles retrieves multiple JSON files from S3
func (c *Client) GetFiles() (map[string]io.Reader, error) {
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

	output := make(map[string]io.Reader)
	for _, object := range objects {
		file, err := getFile(c.s3, *object.Key)
		if err != nil {
			return nil, fmt.Errorf("error getting object %s: %s", *object.Key, err.Error())
		}
		output[strconv.Itoa(year)+"-"+strconv.Itoa(month)+*object.Key] = file
	}

	return output, nil
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
