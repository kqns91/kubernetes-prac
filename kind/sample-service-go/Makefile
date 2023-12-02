.PHONY: build
build:
	docker build . -t sample-service-go:latest
	kind load --name local-dev docker-image sample-service-go:latest