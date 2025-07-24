# Development Guide

This guide helps developers get started with the BOGOWI Blockchain API Go project.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+ installed
- Git configured
- Docker (optional)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/KODESL/bogowi-blockchain-api.git
cd bogowi-blockchain-api

# Set up development environment
make dev-setup

# Copy environment file and configure
cp .env.example .env
# Edit .env with your configuration

# Run the API
make run
```

## 🛠️ Development Workflow

### Code Formatting and Quality

The project uses automated code formatting and quality checks:

```bash
# Format code
make format

# Check formatting
make format-check

# Run linter
make lint

# Run security scan
make security

# Run all pre-commit checks
make pre-commit
```

### Building

```bash
# Build for current platform
make build

# Build for Linux (deployment)
make build-linux

# Build for all platforms
make build-all
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
# Opens coverage.html in browser
```

### Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run
```

## 📋 Available Make Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build` | Build the binary |
| `make test` | Run tests |
| `make format` | Format Go code |
| `make lint` | Run linter |
| `make security` | Run security scan |
| `make pre-commit` | Run all pre-commit checks |
| `make ci` | Run CI checks |
| `make clean` | Clean build artifacts |
| `make run` | Run the application |

## 🔄 Git Workflow

This project uses Git Flow:

### Feature Development

```bash
# Start a new feature
git flow feature start feature-name

# Work on your feature...

# Finish the feature
git flow feature finish feature-name
```

### Pre-commit Hooks

The project has pre-commit hooks that automatically:
- Format your code with `gofmt`
- Run `go vet` for static analysis
- Ensure code quality before commits

### Code Quality Standards

- **Formatting**: All code must be formatted with `gofmt -s`
- **Linting**: Must pass `golangci-lint` checks
- **Testing**: New features should include tests
- **Security**: Code is scanned with `gosec`

## 📁 Project Structure

```
bogowi-blockchain-api/
├── main.go                 # Application entry point
├── internal/               # Private application code
│   ├── api/               # HTTP handlers and routing
│   ├── config/            # Configuration management  
│   └── sdk/               # Blockchain SDK
├── docs/                  # API documentation
├── .github/workflows/     # CI/CD workflows
├── Makefile              # Development commands
├── Dockerfile            # Container configuration
├── go.mod                # Go module definition
└── README.md             # Project documentation
```

## ⚙️ Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Required variables
NETWORK=CAMINO_TESTNET
RPC_URL=https://columbus.camino.network/ext/bc/C/rpc
PRIVATE_KEY=your-private-key

# Contract addresses
BOGO_TOKEN_V2_ADDRESS=0x...
CONSERVATION_NFT_ADDRESS=0x...
# ... other contract addresses

# Optional variables
API_PORT=3001
LOG_LEVEL=info
```

### Local Development

For local development, you can use test/dummy values:

```bash
ENVIRONMENT=development
PRIVATE_KEY=0x0000000000000000000000000000000000000000000000000000000000000000
```

## 🚀 Deployment

The project includes automated deployment via GitHub Actions:

### Manual Deployment

```bash
# Build for Linux
make build-linux

# Copy to server (example)
scp bogowi-api-linux user@server:/path/to/deployment/
```

### GitHub Actions

- **CI**: Runs on every push/PR to `main`/`develop`
- **Deployment**: Automatically deploys `main` branch to production
- **Release**: Creates releases with multi-platform binaries

## 🐛 Debugging

### Local Development

```bash
# Run with verbose logging
LOG_LEVEL=debug make run

# Run tests with verbose output
go test -v ./...
```

### Docker Debugging

```bash
# Build and run in debug mode
docker build -t bogowi-api-debug .
docker run -p 3001:3001 --env-file .env bogowi-api-debug

# View logs
docker logs container-name
```

## 📚 API Documentation

- **Interactive Docs**: Visit `http://localhost:3001/docs` when running
- **OpenAPI Spec**: See `openapi.yaml` for complete API specification
- **Health Check**: `GET /api/health` for service status

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git flow feature start amazing-feature`)
3. **Make** your changes
4. **Run** pre-commit checks (`make pre-commit`)
5. **Commit** your changes
6. **Push** to your fork
7. **Create** a Pull Request

### Code Review Checklist

- [ ] Code is formatted (`make format-check`)
- [ ] Linting passes (`make lint`)
- [ ] Tests pass (`make test`)
- [ ] Security scan passes (`make security`)
- [ ] Documentation updated if needed

## 🆘 Troubleshooting

### Common Issues

**Module not found errors:**
```bash
make deps
```

**Formatting errors:**
```bash  
make format
```

**Build errors:**
```bash
make clean
make build
```

**Port already in use:**
```bash
# Change API_PORT in .env
API_PORT=3002
```

### Getting Help

- Check the [README](README.md) for general information
- Review [GitHub Issues](https://github.com/KODESL/bogowi-blockchain-api/issues)
- Contact the development team

---

Happy coding! 🌊🌍
