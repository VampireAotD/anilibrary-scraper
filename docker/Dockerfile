# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o /bin/scraper ./cmd/app

FROM alpine:latest as final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add tzdata

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    scraper
USER scraper

COPY --from=builder /bin/scraper /bin/

ARG TIMEZONE
ENV TZ=$TIMEZONE

ENTRYPOINT ["/bin/scraper"]