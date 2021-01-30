package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fijixxx/sublog-function/delete"
)

func main() {
	lambda.Start(delete.HandleRequest)
}
