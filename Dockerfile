FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache git ca-certificates make gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/main ./cmd/api

FROM alpine:3.18 AS final

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/app/main"]
