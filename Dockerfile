# Build stage
FROM golang:alpine AS builder

# Install git
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /app/main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy all YAML and YML configuration files
COPY config/config.yaml /app/config/config.yaml
COPY config/configProduction.yaml /app/config/configProduction.yaml


EXPOSE 8080

ENTRYPOINT ["/app/main"]