FROM golang:1.23.3-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY 
