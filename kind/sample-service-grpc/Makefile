.PHONY: build
build:
	docker build . -t sample-service-grpc:latest
	kind load --name local-dev docker-image sample-service-grpc:latest