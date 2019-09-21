#!/bin/bash
# go build -ldflags "-X main.HANDLER=SAVE" -o lambdasave
# zip quehook-save.zip lambdasave
# aws lambda update-function-code --function-name quehook-save --zip-file fileb://quehook-save.zip --region us-east-1
#
# go build -ldflags "-X main.HANDLER=LOAD" -o lambdaload
# zip quehook-load.zip lambdaload
# aws lambda update-function-code --function-name quehook-load --zip-file fileb://quehook-load.zip --region us-east-1
#
# go build -ldflags "-X main.HANDLER=BACKFILL" -o lambdabackfill
# zip quehook-backfill.zip lambdabackfill
# aws lambda update-function-code --function-name quehook-backfill --zip-file fileb://quehook-backfill.zip --region us-east-1
