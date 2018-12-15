FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/snaps/
WORKDIR $GOPATH/src/snaps/
RUN go get -d -v
RUN go build -o /go/bin/main
ENTRYPOINT ["/go/bin/main"]
