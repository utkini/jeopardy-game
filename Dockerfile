# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o jeopardy-game ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/jeopardy-game .
# Copy config files
COPY --from=builder /app/configs ./configs
# Copy templates
COPY --from=builder /app/templates ./templates

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./jeopardy-game"] 