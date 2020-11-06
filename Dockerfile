FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/go

RUN mkdir -p /src/build
COPY . /src/build
WORKDIR /src/build

RUN go build -o . ./...

FROM alpine:3
COPY --from=builder /src/build/website website
COPY ./assets ./assets
COPY ./data ./data
COPY ./views ./views

EXPOSE 8675
ENTRYPOINT ["/website"]