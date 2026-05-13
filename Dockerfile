# =========================
# 1. Build stage
# =========================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auth-service ./cmd/main.go


# =========================
# 2. Runtime stage
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service

EXPOSE 50051

ENTRYPOINT ["/app/auth-service"]