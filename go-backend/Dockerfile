FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /build

# Copy source code first (needed for go mod tidy to work)
COPY . .

# Generate go.sum and download dependencies
RUN go mod tidy
RUN go mod download
RUN go mod verify

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o timelith-backend ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/timelith-backend .

# Create necessary directories
RUN mkdir -p uploads sessions logs

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./timelith-backend"]
