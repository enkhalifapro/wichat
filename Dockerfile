# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# install git
RUN apk add --update --no-cache git
############ Copy the local package files to the container's workspace.
ADD . /go/src/wichat
WORKDIR /go/src/wichat
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build
ENTRYPOINT ./wichat