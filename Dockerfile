FROM golang:1.14.6-alpine3.12

WORKDIR /go/src/app

# Copy App files
COPY app/cmd .

# Wait for Postgres to be up before launching API
COPY app/scripts/wait-for-it.sh .
RUN chmod +x ./wait-for-it.sh

# Install GIT package
RUN apk update && apk upgrade && \
    apk add bash git openssh

# Install GoLang and compile
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v ./...
