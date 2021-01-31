FROM golang:latest as builder
ENV GOPATH=
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dist/main upsert/cmd/lambda/main.go && mv entry.sh /dist/

FROM alpine:3.12.3
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
COPY --from=builder /dist/ /
RUN chmod 755 /usr/bin/aws-lambda-rie && chmod 755 /entry.sh
ENTRYPOINT [ "/entry.sh" ]
CMD [ "/main" ]