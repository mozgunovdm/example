FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod tidy -compat=1.17

ENTRYPOINT go run cmd/employe/main.go
