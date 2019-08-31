package handlers

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tidwall/gjson"

	"github.com/forstmeier/comana/storage"
)

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
	log.Printf("save request: %s", req.Body)

	if req.Source == "" || (req.Source != "aws.events" && req.Source != "comana.backfill") {
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "source must be cloudwatch event or backfill",
			IsBase64Encoded: false,
		}, errors.New("source must be cloudwatch event or backfill")
	}

	year, month, day, hour := req.Year, req.Month, req.Day, req.Hour
	if req.Source == "aws.events" {
		current := time.Now().Add(time.Hour * -5) // 1 hour prior + 4 hour UTC-EST difference
		year, _, day = current.Date()
		month = int(current.Month())
		hour = current.Hour()
	}
	log.Printf("source: %s, year: %d, month: %d, start day: %d, end day: %d", req.Source, year, month, day, hour)

	url := fmt.Sprintf("https://data.gharchive.org/%d-%02d-%02d-%d.json.gz", year, month, day, hour)
	log.Printf("gh archive url: %s", url)

	file, err := download(url)
	if err != nil {
		log.Println("error retrieving archive file: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error retrieving archive file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	scanner, err := unzip(file)
	if err != nil {
		log.Println("error unzipping archive file: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error unzipping archive file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	reader, err := parse(scanner)
	if err != nil {
		log.Println("error parsing archive file: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error parsing archive file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	if err := s.PutFile(year, month, day, hour, reader); err != nil {
		log.Println("error saving report file: " + err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			Body:            "error saving report file: " + err.Error(),
			IsBase64Encoded: false,
		}, err
	}

	log.Println("successful save")
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}
