FROM golang:1.12.1-alpine3.9 as build
RUN apk add --update git

ENV CGO_ENABLED=0
ARG GOPROXY

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN go test ./...
RUN go install . ./cmd/...

FROM alpine:3.9
#entrypoint ...
COPY --from=build /go/bin/cmd /hafenhause
