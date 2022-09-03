## build stage
FROM golang:1.19-alpine AS build-env
RUN apk --no-cache add build-base git curl tzdata
ADD . /build
WORKDIR /build/cmd/app

RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -o=scraper

## final stage
FROM alpine:latest
COPY --from=build-env /build/cmd/app /cmd/app

WORKDIR /cmd/app

ENTRYPOINT ["./scraper"]