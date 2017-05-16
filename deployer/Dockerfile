FROM golang:1.8.0-alpine
RUN apk update && apk add docker
ADD main.go main.go
CMD go run main.go