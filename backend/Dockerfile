# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/api ./cmd/api

FROM alpine:3.20
WORKDIR /app

ENV APP_ENV=production \
    HTTP_PORT=8080 \
    TENANT_HEADER=X-Tenant-ID \
    LOG_LEVEL=info

COPY --from=builder /tmp/api /usr/local/bin/api
EXPOSE 8080

CMD ["api"]
