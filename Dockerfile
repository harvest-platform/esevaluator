FROM golang:latest

RUN mkdir -p /go/src/github.com/harvest-platform/esevaluator
COPY . /go/src/github.com/harvest-platform/esevaluator

WORKDIR /go/src/github.com/harvest-platform/esevaluator

RUN go build -o esevaluator ./cmd/esevaluator

CMD ["./esevaluator"]
