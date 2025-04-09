# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.1

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -ldflags '-w -s -extldflags "-static"' -o /bin/scraper ./cmd/app

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add tzdata

ARG UID=10001
ARG USER=scraper

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    $USER

USER $USER

COPY --from=builder /bin/scraper /bin/

ARG TIMEZONE
ENV TZ=$TIMEZONE

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/bin/scraper"]