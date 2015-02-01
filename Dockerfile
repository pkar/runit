FROM golang

WORKDIR /go
ADD . /go/src/github.com/pkar/runit
RUN go get ./...
RUN go install github.com/pkar/runit/cmd/runit
