# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application and migrator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/notify-srv ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/migrator ./cmd/migrator/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates netcat-openbsd

# Copy the built binaries from the builder stage
COPY --from=builder /app/notify-srv .
COPY --from=builder /app/migrator .
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8098

ENTRYPOINT ["./entrypoint.sh"]