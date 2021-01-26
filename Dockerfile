FROM public.ecr.aws/lambda/go:latest
WORKDIR /var/task
COPY sublog-function .
CMD [ "main" ]