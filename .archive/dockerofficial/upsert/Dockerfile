FROM golang:latest as builder
ENV GOPATH=
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /dist/main upsert/cmd/lambda/main.go

FROM public.ecr.aws/lambda/provided:al2
COPY bootstrap ${LAMBDA_RUNTIME_DIR}
RUN chmod 755 ${LAMBDA_RUNTIME_DIR}/bootstrap
COPY --from=builder /dist/ /
CMD [ "/main" ]