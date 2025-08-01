FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

# Create final image
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-shortener .

# Copy static assets
COPY internal/web/templates ./templates
COPY internal/web/styles.css ./styles.css
COPY internal/web/main.js ./js

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./url-shortener"]
