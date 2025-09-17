# syntax=docker/dockerfile:1

FROM golang:1.24-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache git

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY ./cmd/api ./cmd/api

# Build the binary
RUN go build -o /teams-api ./cmd/api

# Expose port
EXPOSE 8080

# Run binary
CMD ["/teams-api"]
