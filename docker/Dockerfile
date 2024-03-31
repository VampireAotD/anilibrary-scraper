## modules stage
FROM golang:1.22-alpine AS modules
ADD go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

## build stage
FROM golang:1.22-alpine AS builder
RUN apk --no-cache add tzdata

COPY --from=modules /go/pkg /go/pkg

ADD . /build
WORKDIR /build/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o=scraper

## final stage
FROM scratch

COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=builder /build/cmd/app /cmd/bin

ARG TIMEZONE

ENV TZ=$TIMEZONE

WORKDIR /cmd/bin

ENTRYPOINT ["./scraper"]