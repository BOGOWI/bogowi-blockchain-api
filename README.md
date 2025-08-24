# BOGOWI Blockchain API

> A high-performance Go REST API for the BOGOWI blockchain ecosystem - enabling conservation initiatives and sustainable travel on the Camino network.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-green.svg)](https://gin-gonic.com/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

## 🌊 Overview

The BOGOWI Blockchain API is a production-ready Go backend service that provides RESTful endpoints for interacting with the BOGOWI ecosystem smart contracts on the Camino blockchain. It enables gamified conservation efforts, NFT management, and DAO governance operations.

### Key Features

- **🚀 High Performance**: Built with Go and Gin framework for optimal throughput
- **🔐 Enterprise Security**: JWT authentication, rate limiting, and CORS protection
- **📚 API Documentation**: Auto-generated OpenAPI 3.0 specification with Redoc UI
- **🔄 Smart Contract Integration**: Direct interaction with BOGO tokens, NFTs, and DAO contracts  
- **📊 Real-time Data**: Live blockchain state queries and transaction monitoring
- **🐳 Docker Ready**: Containerized deployment with health checks
- **🎯 Production Tested**: Running on AWS with nginx reverse proxy

## 🏗️ Architecture

```
api/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication, CORS, rate limiting
│   ├── models/          # Data structures and DTOs
│   ├── services/        # Business logic layer
│   └── blockchain/      # Smart contract interactions
├── pkg/
│   └── sdk/            # BOGOWI blockchain SDK
├── docs/               # API documentation assets
├── docker/             # Docker configuration
└── deployments/        # Kubernetes/deployment configs
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional)
- Access to Camino network RPC endpoint

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/KODESL/bogowi-blockchain-api.git
   cd bogowi-blockchain-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the API server**
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:3001`

### Docker Deployment

```bash
# Build and run with Docker
docker build -t bogowi-api .
docker run -p 3001:3001 --env-file .env bogowi-api

# Or use Docker Compose
docker-compose up -d
```

## 📘 API Documentation

### Interactive Documentation

Visit `http://localhost:3001/docs` to explore the interactive API documentation powered by Redoc.

### Core Endpoints

#### Health & Status
- `GET /api/health` - Service health check with contract addresses
- `GET /api/contracts` - Current smart contract addresses and ABIs

#### Token Operations  
- `GET /api/tokens/bogo/balance/{address}` - Get BOGO token balance
- `GET /api/tokens/flavored/{type}/balance/{address}` - Get flavored token balance (OCEAN, EARTH, WILDLIFE)
- `POST /api/tokens/exchange` - Exchange flavored tokens for BOGO

#### NFT Management
- `GET /api/nfts/{address}/portfolio` - Get user's NFT portfolio
- `GET /api/nfts/collection/{id}` - Get NFT collection details
- `POST /api/nfts/mint` - Mint NFTs (authorized users only)

#### User Management
- `GET /api/users/{address}/profile` - Get user profile and stats
- `GET /api/users/{address}/transactions` - Get transaction history
- `POST /api/users/register` - Register new user profile

#### DAO Operations
- `GET /api/dao/proposals` - List active governance proposals
- `GET /api/dao/treasury` - Get treasury balance and info
- `POST /api/dao/vote` - Cast vote on proposal (authenticated)

#### Analytics & Reporting
- `GET /api/analytics/conservation` - Conservation impact metrics
- `GET /api/analytics/tokens` - Token distribution statistics
- `GET /api/leaderboard` - Community leaderboard

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Rate Limiting

API requests are rate-limited to prevent abuse:
- **Public endpoints**: 100 requests/hour per IP
- **Authenticated endpoints**: 1000 requests/hour per user

## 🔧 Development

### Project Structure

```
bogowi-blockchain-go/
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routing
│   │   ├── router.go      # Route definitions
│   │   ├── system.go      # System endpoints
│   │   ├── tokens.go      # Token endpoints
│   │   └── handlers.go    # NFT/Rewards/DAO handlers
│   ├── config/            # Configuration management
│   │   └── config.go      # Config loading and AWS SSM
│   └── sdk/               # Blockchain interaction layer
│       ├── sdk.go         # Main SDK implementation
│       └── abi.go         # Contract ABIs
├── Dockerfile             # Container configuration
├── .env.template          # Environment template
└── README.md
```

### Building

```bash
# Development build
go build -o bogowi-api

# Production build (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bogowi-api
```

### Testing

```bash
# Run tests
go test ./...

# Test with coverage
go test -cover ./...
```

## 🚀 Deployment


## 🚀 Deployment

### Docker Deployment

```bash
# Build and run with Docker
docker build -t bogowi-api .
docker run -p 3001:3001 --env-file .env bogowi-api

# Or use Docker Compose
docker-compose up -d
```

### Environment Configuration

Copy the example environment file and configure:
```bash
cp .env.example .env
# Edit .env with your configuration
```

### Production Considerations

- Use HTTPS in production
- Configure proper rate limiting
- Set up monitoring and logging
- Use environment variables for sensitive data
- Implement proper backup strategies

## 🔗 Related Repositories

- [Smart Contracts](https://github.com/BOGOWI/bogowi-contracts) - BOGOWI ecosystem smart contracts
- [Frontend](https://github.com/BOGOWI/bogowi-frontend) - Web application frontend

## 🌍 Community

Join our conservation community:
- Website: [bogowi.com](https://bogowi.com)
- Discord: [Join our community](https://discord.gg/bogowi)
- Twitter: [@BOGOWI](https://twitter.com/BOGOWI)

---

**Building a sustainable future through blockchain technology** 🌱
