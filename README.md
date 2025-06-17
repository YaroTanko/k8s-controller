# Kubernetes Controller

A Kubernetes controller application that manages custom resources and performs operations based on Kubernetes events.

[![CI/CD](https://github.com/YaroTanko/k8s-controller/actions/workflows/ci.yml/badge.svg)](https://github.com/YaroTanko/k8s-controller/actions/workflows/ci.yml)

## Features

- Kubernetes controller functionality for custom resource management
- HTTP server with health checks and API endpoints
- Structured logging with configurable levels
- Request logging middleware with detailed metrics
- CLI interface with subcommands using Cobra
- Configuration via environment variables or command-line flags

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Kubernetes cluster (for controller functionality)
- Docker (for containerized deployment)

### Installation

#### From Source

```bash
# Clone the repository
git clone https://github.com/YaroTanko/k8s-controller.git
cd k8s-controller

# Build the binary
make build

# Run the controller
./bin/k8s-controller
```

#### Using Docker

```bash
# Build the Docker image
make docker-build

# Run the container
docker run -p 8080:8080 k8s-controller:latest
```

## Usage

### Command Line Interface

The application provides a CLI with several subcommands:

```bash
# Show help and available commands
./bin/k8s-controller --help

# Run the Kubernetes controller
./bin/k8s-controller serve --namespace default

# Start the HTTP server
./bin/k8s-controller server --port 9090

# Show version information
./bin/k8s-controller version
```

### Global Flags

- `--log-level`, `-l`: Set logging level (trace, debug, info, warn, error)
- `--kubeconfig`, `-k`: Path to kubeconfig file
- `--namespace`, `-n`: Kubernetes namespace to operate in

### Server Mode

```bash
./bin/k8s-controller server [flags]

Flags:
  --debug           Enable debug mode with detailed request logging
  --port int        HTTP server port (default 8080)
```

### Controller Mode

```bash
./bin/k8s-controller serve [flags]

Flags:
  --leader-elect    Enable leader election
  --workers int     Number of worker threads (default 2)
```

## Configuration

Configuration can be provided via environment variables or command-line flags:

| Environment Variable | Flag | Description | Default |
|----------------------|------|-------------|---------|
| K8S_CONTROLLER_LOG_LEVEL | --log-level | Logging level | info |
| K8S_CONTROLLER_KUBECONFIG | --kubeconfig | Path to kubeconfig | |
| K8S_CONTROLLER_NAMESPACE | --namespace | Kubernetes namespace | |
| K8S_CONTROLLER_SERVER_PORT | --port | HTTP server port | 8080 |

## Development

### Build and Test

```bash
# Build the binary
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Clean build artifacts
make clean
```

### Docker Development

```bash
# Build Docker image
make docker-build

# Run Docker container
docker run -p 8080:8080 k8s-controller:latest
```

### Building for Multiple Platforms

```bash
make build-all
```

## Deployment

The application is packaged as a distroless container for secure deployment in Kubernetes.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-controller
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-controller
  template:
    metadata:
      labels:
        app: k8s-controller
    spec:
      containers:
      - name: k8s-controller
        image: ghcr.io/yarotanko/k8s-controller:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: K8S_CONTROLLER_LOG_LEVEL
          value: "info"
```

## Project Structure

```
.
├── cmd/                # Command line interface
│   ├── root.go         # Root command and global flags
│   ├── serve.go        # Kubernetes controller command
│   ├── server.go       # HTTP server command
│   └── version.go      # Version information command
├── pkg/                # Core packages
│   ├── config/         # Configuration handling
│   ├── logger/         # Structured logging
│   └── middleware/     # HTTP middleware components
├── Dockerfile          # Distroless container definition
├── Makefile            # Build and development tasks
└── main.go             # Application entry point
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/amazing-feature`)
3. Commit your Changes (`git commit -m 'Add some amazing feature'`)
4. Push to the Branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request