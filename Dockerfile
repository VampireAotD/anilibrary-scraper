## build stage
FROM golang:1.18.3-alpine AS build-env
RUN apk --no-cache add build-base git curl
ADD . /build
WORKDIR /build/app/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o=scrapper

## final stage
FROM alpine:latest
COPY --from=build-env /build/app/cmd/app /app/cmd/app
COPY --from=build-env /build/.env .env

WORKDIR /app/cmd/app

ENTRYPOINT ["./scrapper"]