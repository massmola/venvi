# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 is used because we use modernc.org/sqlite (pure go)
RUN CGO_ENABLED=0 GOOS=linux go build -o venvi main.go

# Run stage
FROM alpine:latest

RUN apk add --no-cache \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/venvi /app/venvi

# Copy static assets and migrations
COPY --from=builder /app/pb_public /app/pb_public
COPY --from=builder /app/pb_migrations /app/pb_migrations

# Create data directory for SQLite
RUN mkdir /app/pb_data

# Expose the default PocketBase port
EXPOSE 8080

# Command to run the application
# Use --dir to specify the data directory and --http for the port
CMD ["/app/venvi", "serve", "--http=0.0.0.0:8080", "--dir=/app/pb_data"]
