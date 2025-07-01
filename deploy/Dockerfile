# -------- STAGE 1: Build --------
FROM golang:1.24.3-alpine AS builder

# Install git (some Go modules need it)
RUN apk add --no-cache git

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/app

# -------- STAGE 2: Run --------
FROM alpine:latest

# Install certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
