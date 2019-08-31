#!/bin/bash
go build -ldflags "-X main.HANDLER=SAVE" -o lambdasave
zip comana-save.zip lambdasave
aws lambda update-function-code --function-name comana-save --zip-file fileb://comana-save.zip --region us-east-1

go build -ldflags "-X main.HANDLER=LOAD" -o lambdaload
zip comana-load.zip lambdaload
aws lambda update-function-code --function-name comana-load --zip-file fileb://comana-load.zip --region us-east-1

go build -ldflags "-X main.HANDLER=BACKFILL" -o lambdabackfill
zip comana-backfill.zip lambdabackfill
aws lambda update-function-code --function-name comana-backfill --zip-file fileb://comana-backfill.zip --region us-east-1
