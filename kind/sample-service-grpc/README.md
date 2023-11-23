```sh
docker build . -t sample-service-grpc

kind load --name local-dev docker-image sample-service-grpc:latest
```