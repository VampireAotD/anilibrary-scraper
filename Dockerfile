## modules stage
FROM golang:1.19-alpine AS modules
ADD go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

## build stage
FROM golang:1.19-alpine AS builder
RUN apk --no-cache add tzdata

COPY --from=modules /go/pkg /go/pkg

ADD . /build
WORKDIR /build/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o=scraper

## final stage
FROM alpine:latest

COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=builder /build/cmd/app /cmd/bin

WORKDIR /cmd/bin

ENTRYPOINT ["./scraper"]