FROM golang:1.23

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN usermod -u 1000 www-data

RUN chmod -R 777 /app