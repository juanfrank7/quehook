#!/bin/bash

go build -ldflags "-X main.HANDLER=CREATE" -o querycreate
zip quehook-querycreate.zip querycreate
aws lambda update-function-code --function-name quehook-query-create --zip-file fileb://quehook-querycreate.zip --region us-east-1

go build -ldflags "-X main.HANDLER=RUN" -o queryrun
zip quehook-queryrun.zip queryrun
aws lambda update-function-code --function-name quehook-query-run --zip-file fileb://quehook-queryrun.zip --region us-east-1

go build -ldflags "-X main.HANDLER=DELETE" -o querydelete
zip quehook-querydelete.zip querydelete
aws lambda update-function-code --function-name quehook-query-delete --zip-file fileb://quehook-querydelete.zip --region us-east-1

go build -ldflags "-X main.HANDLER=SUBSCRIBE" -o querysubscribe
zip quehook-querysubscribe.zip querysubscribe
aws lambda update-function-code --function-name quehook-subscription-subscribe --zip-file fileb://quehook-querysubscribe.zip --region us-east-1

go build -ldflags "-X main.HANDLER=UNSUBSCRIBE" -o queryunsubscribe
zip quehook-queryunsubscribe.zip queryunsubscribe
aws lambda update-function-code --function-name quehook-subscription-unsubscribe --zip-file fileb://quehook-queryunsubscribe.zip --region us-east-1
