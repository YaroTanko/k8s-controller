# Build Stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/k8s-controller main.go

# Final Stage with Distroless
FROM gcr.io/distroless/static:nonroot

# Copy the binary from builder
COPY --from=builder /app/bin/k8s-controller /k8s-controller

# Use nonroot user
USER nonroot:nonroot

# Expose default server port
EXPOSE 8080

# Set command
ENTRYPOINT ["/k8s-controller"]
CMD ["serve"]