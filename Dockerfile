# Stage 1: Build
FROM golang:alpine AS builder

# Install necessary packages
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .
    
# Build the Go app
RUN go build -o /app/main ./main.go

# Final Stage: Runtime
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /main

# Command to run the binary
CMD ["/main"]
