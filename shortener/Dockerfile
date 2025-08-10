FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/shortener

# Create final image
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-shortener .

# Copy web assets
COPY web/ ./web/

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./url-shortener"]
