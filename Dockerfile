# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Optional: static binary biar ringan & aman
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Install CA certificates (penting kalau ada HTTPS request)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

EXPOSE 8081

CMD ["./app"]