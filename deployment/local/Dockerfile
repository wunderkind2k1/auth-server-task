# Build stage
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy the entire server directory
COPY server/ ./

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-server

# Final stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/auth-server .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./auth-server"]
