# Base
FROM --platform=$BUILDPLATFORM golang:1.23.0-alpine AS builder
RUN apk add --no-cache build-base
WORKDIR /app
COPY . /app
ARG TARGETARCH
RUN go mod download
RUN GOOS=linux GOARCH=$TARGETARCH go build -o tennis-tg-bot .

# Release
FROM alpine:latest
RUN apk -U upgrade --no-cache \
	&& apk add --no-cache bind-tools ca-certificates
COPY --from=builder /app/tennis-tg-bot /usr/local/bin/

ARG GIT_SHA

LABEL org.opencontainers.image.authors="@tubopo" \
	org.opencontainers.image.description="Telegram tennis training scheduler bot" \
	org.opencontainers.image.licenses="MIT" \
	org.opencontainers.image.source="https://github.com/tubopo/tennis-tg-bot" \
	org.opencontainers.image.title="tennis-tg-bot" \
	org.opencontainers.image.revision="${GIT_SHA}"

ENTRYPOINT ["tennis-tg-bot"]
