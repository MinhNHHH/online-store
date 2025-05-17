# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy .env file
COPY .env .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o store cmd/store.go

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/store .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy .env file to final stage
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["./store"] 