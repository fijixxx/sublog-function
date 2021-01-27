package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fijixxx/sublog-function/upsert"
)

func main() {
	lambda.Start(upsert.HandleRequest)
}
