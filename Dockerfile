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

ARG GIT_SHA

LABEL org.opencontainers.image.authors="@tubopo" \
	org.opencontainers.image.description="Telegram tennis traning scheduler bot" \
	org.opencontainers.image.licenses="MIT" \
	org.opencontainers.image.source="https://github.com/tubopo/tennis-tg-bot" \
	org.opencontainers.image.title="tennis-tg-bot" \
	org.opencontainers.image.revision="${GIT_SHA}"

ENTRYPOINT ["tennis-tg-bot"]
