# =========================
# 1. Build Stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 빌드 (entrypoint: cmd/main.go)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/main.go


# =========================
# 2. Runtime Stage
# =========================
FROM alpine:3.19

WORKDIR /app

RUN mkdir -p /app/log
RUN mkdir -p /certs

COPY --from=builder /app/app .

COPY cmd/app.env ./app.env
#COPY certs/ /certs/

EXPOSE 9083

# 실행
ENTRYPOINT ["./app"]
