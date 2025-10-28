# Use official Go image as base
FROM golang:1.21-alpine

# Set working directory inside container
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o user-service ./cmd/user-service/main.go

# Expose port (if your service listens on 8080)
EXPOSE 8080

# Run binary
CMD ["./user-service"]