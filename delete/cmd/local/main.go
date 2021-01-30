package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/fijixxx/sublog-function/delete"
)

func main() {
	ctx := new(context.Context)
	e := new(events.S3Event)
	delete.HandleRequest(*ctx, *e)
}
