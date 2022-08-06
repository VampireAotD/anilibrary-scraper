## build stage
FROM golang:1.18.3-alpine AS build-env
RUN apk --no-cache add build-base git curl tzdata
ADD . /build
WORKDIR /build/app/cmd/app

RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o=scraper

## final stage
FROM alpine:latest
COPY --from=build-env /build/app/cmd/app /app/cmd/app
COPY --from=build-env /build/.env .env

WORKDIR /app/cmd/app

ENTRYPOINT ["./scraper"]