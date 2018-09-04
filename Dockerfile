####
FROM golang:1.11-alpine3.8 AS build

# go workspace
RUN mkdir -p /go/bin/
RUN mkdir -p /go/pkg/
RUN mkdir -p /go/src/WebApi_Go/

#installing packages
RUN apk update
RUN apk upgrade
RUN apk add git
RUN go get -u github.com/golang/dep/cmd/dep

#build
WORKDIR /go/src/WebApi_Go/
COPY . .
RUN dep init
RUN dep ensure
RUN go build

FROM alpine:3.8 AS prod
COPY --from=build /go/src/WebApi_Go/WebApi_Go .
CMD ["./WebApi_Go"]