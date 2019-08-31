package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/tidwall/gjson"
)

// Invoker wraps logic for triggering a Lambda
type Invoker interface {
	Invoke(payload []byte) (int64, string, error)
}

type client struct {
	lambda *lambda.Lambda
}

func (i *client) Invoke(payload []byte) (int64, string, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String("comana-save"),
		Payload:      payload,
	}

	result, err := i.lambda.Invoke(input)
	return *result.StatusCode, string(result.Payload), err
}

// NewInvoke generates an Invoke implementation with an active client
func NewInvoke() Invoker {
	return &client{
		lambda: lambda.New(session.New()),
	}
}

// BackfillData pulls in historic data for stat updates
func BackfillData(req Request, client Invoker) (events.APIGatewayProxyResponse, error) {
	log.Printf("backfill request: %s", req.Body)

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
	log.Printf("year: %d, month: %d, start day: %d, end day: %d", year, month, startDay, endDay)

	finished := make(chan bool, 1)
	errs := make(chan error, 1)
	var wg sync.WaitGroup

	for day := startDay; day <= endDay; day++ {
		for hour := 0; hour < 24; hour++ {
			url := fmt.Sprintf("%s/%d-%02d-%02d-%d.json.gz", "https://data.gharchive.org", year, month, day, hour)
			log.Printf("gh archive url: %s", url)

			wg.Add(1)
			go func(year, month, day, hour int, url string) {
				defer wg.Done()
				payloadRequest := Request{
					Source: "comana.backfill",
					Year:   year,
					Month:  month,
					Day:    day,
					Hour:   hour,
				}

				payload, err := json.Marshal(payloadRequest)
				if err != nil {
					errs <- fmt.Errorf("payload marshalling error for %s: %s", url, err.Error())
				}

				code, resp, err := client.Invoke(payload)
				if err != nil {
					errs <- fmt.Errorf("lambda invocation error for %s: %s", url, err.Error())
				}

				log.Printf("save lambda status code: %d, response: %s", code, resp)
			}(int(year), int(month), int(day), hour, url)
		}
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	select {
	case <-finished:
	// NOTE: possibly expand to collect all goroutine errors
	case err := <-errs:
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:      500,
				Body:            "error invoking lambda: " + err.Error(),
				IsBase64Encoded: false,
			}, err
		}
	}

	log.Println("successful backfill")
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            "success",
		IsBase64Encoded: false,
	}, nil
}
