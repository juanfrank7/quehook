package handlers

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tidwall/gjson"

	"github.com/forstmeier/comana/storage"
)

// Request provides a generalization of CloudWatch and API Gateway events
type Request struct {
	Body       string            `json:"body"`
	HTTPMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers"`
	Source     string            `json:"source"`
}

var download = func(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

var unzip = func(input []byte) (*bufio.Scanner, error) {
	r := strings.NewReader(string(input))

	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(gz), nil
}

var parse = func(s *bufio.Scanner) (io.Reader, error) {
	data := make(map[string]map[string]int)
	for s.Scan() {
		line := s.Text()

		event := gjson.Get(line, "type").String()
		repo := gjson.Get(line, "repo.name").String()

		if _, repoExists := data[repo]; repoExists {
			if _, eventExists := data[repo][event]; eventExists {
				data[repo][event]++
			} else {
				data[repo][event] = 1
			}
		} else {
			data[repo] = map[string]int{
				event: 1,
			}
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(b)), nil
}

// SaveData pulls in and parses GitHub Archive data
func SaveData(req Request, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	if req.Source == "" || req.Source != "aws.events" {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "request source must be cloudwatch event",
			IsBase64Encoded: false,
		}, errors.New("request source must be cloudwatch event")
	}

	current := time.Now().Add(time.Hour * -5) // 1 hour prior + 4 hour UTC-EST difference
	year, m, day := current.Date()
	month := int(m)
	hour := current.Hour()

	url := fmt.Sprintf("%s/%d-%02d-%02d-%d.json.gz", "https://data.gharchive.org", year, month, day, hour)
	file, err := download(url)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error retrieving archieve file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	scanner, err := unzip(file)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error unzipping archieve file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	reader, err := parse(scanner)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error parsing archieve file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	if err := s.PutFile(reader); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error saving report file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}

// LoadData retrieves and returns GitHub Archive reports
func LoadData(s storage.Storage) (events.APIGatewayProxyResponse, error) {
	paths, err := s.GetPaths()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error loading report filepaths: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	pathsObject := map[string][]string{
		"paths": paths,
	}

	output, err := json.Marshal(pathsObject)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error marshalling output: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"result_count": strconv.Itoa(len(paths)),
		},
		Body:            string(output),
		IsBase64Encoded: false,
	}, err
}

// BackfillData pulls in historic data for stat updates
func BackfillData(req Request, s storage.Storage) (events.APIGatewayProxyResponse, error) {
	if req.Headers["COMANA_SECRET"] != os.Getenv("COMANA_SECRET") {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "incorrect secret received: " + req.Headers["COMANA_SECRET"],
			IsBase64Encoded: false,
		}, errors.New("incorrect secret received: " + req.Headers["COMANA_SECRET"])
	}

	year := gjson.Get(req.Body, "year").Int()
	month := gjson.Get(req.Body, "month").Int()
	startDay := gjson.Get(req.Body, "start_day").Int()
	endDay := gjson.Get(req.Body, "end_day").Int()

	for i := startDay; i <= endDay; i++ {
		for j := 0; j < 24; j++ {
			url := fmt.Sprintf("%s/%d-%02d-%02d-%d.json.gz", "https://data.gharchive.org", year, month, i, j)
			file, err := download(url)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error retrieving backfill archieve file: " + err.Error(),
					IsBase64Encoded: false,
				}, err
			}

			scanner, err := unzip(file)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error unzipping backfill archieve file: " + err.Error(),
					IsBase64Encoded: false,
				}, err
			}

			reader, err := parse(scanner)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error parsing backfill archieve file: " + err.Error(),
					IsBase64Encoded: false,
				}, err
			}

			if err := s.PutFile(reader); err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode:      500,
					Body:            "error saving backfill report file: " + err.Error(),
					IsBase64Encoded: false,
				}, err
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}
