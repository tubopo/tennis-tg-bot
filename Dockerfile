# Base
FROM golang:1.23.0-alpine AS builder
RUN apk add --no-cache build-base
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build .

# Release
FROM alpine:3.20.2
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates
COPY --from=builder /app/tennis-tg-bot /usr/local/bin/

ENTRYPOINT ["tennis-tg-bot"]
