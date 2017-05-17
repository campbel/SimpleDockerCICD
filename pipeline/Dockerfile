FROM golang:1.8.0-alpine
RUN apk update && apk add docker
ADD main.go /go/src/pipeline/main.go
ADD pipeline.json pipeline.json
RUN go install pipeline
CMD pipeline