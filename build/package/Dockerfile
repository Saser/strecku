FROM golang:1.12.9

RUN mkdir /strecku
WORKDIR /strecku

COPY go.mod go.sum ./
RUN go mod download

COPY configs ./configs
COPY internal ./internal
