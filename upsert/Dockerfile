FROM golang:latest as builder
ENV GOPATH=
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build upsert/cmd/lambda/main.go

FROM public.ecr.aws/lambda/go:latest
WORKDIR /var/task
COPY --from=builder /go/main .
CMD [ "main" ]