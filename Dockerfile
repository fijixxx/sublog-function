FROM golang:latest as builder
WORKDIR /go/src/github.com/fijixxx/gosuburi/
RUN go get -v -u github.com/BurntSushi/toml \
    github.com/google/uuid \
    github.com/aws/aws-lambda-go/events \
    github.com/aws/aws-lambda-go/lambda \
    github.com/aws/aws-sdk-go/aws \
    github.com/aws/aws-sdk-go/service \
    github.com/guregu/dynamo
COPY main.go .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

FROM public.ecr.aws/lambda/go:latest
WORKDIR /var/task
COPY --from=builder /go/src/github.com/fijixxx/gosuburi/main .
CMD [ "main" ]