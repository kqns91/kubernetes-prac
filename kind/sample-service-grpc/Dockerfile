# syntax=docker/dockerfile:1

# 1. build
FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /api

# 2. deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /api /api

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/api"]
