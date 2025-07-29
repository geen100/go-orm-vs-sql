FROM golang:1.24-alpine AS builder
# Install git and ca-certificates (needed to be able to call HTTPS)
RUN apk --no-cache add ca-certificates git
WORKDIR /app
# Copy go mod files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the binary from builder
COPY --from=builder /app/main .
# Expose port
EXPOSE 8080
# Command to run
CMD ["./main"]

# Development stage (with air for hot reloading)
FROM golang:1.24-alpine AS development
# Install air for hot reloading
RUN go install github.com/cosmtrek/air@latest
# Install dependencies
RUN apk --no-cache add git
WORKDIR /app
# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
# Expose port
EXPOSE 8080
EXPOSE 8081
# Air will handle the hot reloading
CMD ["air", "-c", ".air.toml"]