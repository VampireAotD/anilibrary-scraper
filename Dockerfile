## build stage
FROM golang:1.18.3-alpine AS build-env
RUN apk --no-cache add build-base git curl
ADD . /build
WORKDIR /build/app/cmd/app

RUN GOOS=linux go build -ldflags="-w -s" -o=scrapper

## final stage
FROM alpine:latest
COPY --from=build-env /build/app/cmd/app/ /app/

WORKDIR /app

ENTRYPOINT ["./scrapper"]