FROM golang:1-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/go

RUN mkdir -p /src/build
COPY ./cmd /src/build/cmd
COPY ./internal /src/build/internal
COPY go.mod go.sum /src/build/
WORKDIR /src/build

RUN go build -o . ./...

FROM alpine:3
COPY --from=builder /src/build/website website
COPY ./assets ./assets
COPY ./data ./data
COPY ./views ./views
RUN mkdir ./certs

EXPOSE 8675
CMD ["/website"]