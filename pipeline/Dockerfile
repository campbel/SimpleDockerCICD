FROM golang:1.8.0-alpine
RUN apk update && apk add docker
ADD main.go main.go
ADD pipeline.json pipeline.json
CMD go run main.go