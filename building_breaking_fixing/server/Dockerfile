# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.22.4 AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .

##
## Deploy
##
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --chown=nonroot:nonroot --from=build /app/api /api
EXPOSE 8080/tcp
USER nonroot:nonroot

ENTRYPOINT ["./api"]
