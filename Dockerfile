FROM golang:alpine

WORKDIR /go/src/github.com/wisvch/flitcie-api
ADD . /go/src/github.com/wisvch/flitcie-api
RUN go build -o /go/bin/flitcie-api

ENTRYPOINT ["/go/bin/flitcie-api"]
EXPOSE 8080
