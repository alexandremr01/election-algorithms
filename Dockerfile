FROM golang:1.22.2-bullseye

WORKDIR /app

RUN go install golang.org/x/tools/cmd/goimports@latest

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2
