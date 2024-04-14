## modules stage
FROM golang:1.22-alpine AS modules
WORKDIR /modules
ADD go.mod go.sum ./
RUN go mod download

## build stage
FROM golang:1.22-alpine AS builder

COPY --from=modules /go/pkg /go/pkg

ADD . /build
WORKDIR /build/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o=scraper

## final stage
FROM alpine as final

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

COPY --from=builder /build/cmd/app /bin/

ARG TIMEZONE
ENV TZ=$TIMEZONE

ENTRYPOINT ["/bin/scraper"]