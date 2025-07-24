# Contributing to BOGOWI Blockchain API

Thank you for your interest in contributing to the BOGOWI Blockchain API! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git and Git Flow
- Docker (for containerized development)
- Access to Camino network

### Setup Development Environment

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/bogowi-blockchain-api.git
   cd bogowi-blockchain-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your development configuration
   ```

4. **Run tests to ensure setup**
   ```bash
   go test ./...
   ```

## 🔄 Development Workflow

We use **Git Flow** for branch management:

### Feature Development

1. **Start a new feature**
   ```bash
   git flow feature start feature-name
   ```

2. **Develop your feature**
   - Make your changes
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   go test ./...
   go test -race ./...
   go vet ./...
   ```

4. **Finish the feature**
   ```bash
   git flow feature finish feature-name
   ```

### Bug Fixes & Hotfixes

For critical production issues:

```bash
git flow hotfix start issue-description
# Fix the issue
git flow hotfix finish issue-description
```

## 📝 Code Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Use `go vet` to check for common mistakes
- Follow Go naming conventions

### Code Quality

Before submitting, ensure your code passes:

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run tests
go test ./...

# Test with race detection
go test -race ./...

# Check test coverage
go test -cover ./...
```

### Code Structure

```
internal/
├── api/           # HTTP handlers and routing
├── config/        # Configuration management
├── middleware/    # HTTP middleware
├── models/        # Data structures
├── services/      # Business logic
└── blockchain/    # Smart contract interactions
```

## ✅ Testing Guidelines

### Writing Tests

- Write unit tests for all new functions and methods
- Use table-driven tests where appropriate
- Mock external dependencies
- Test both success and error cases

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/api/

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📚 Documentation

### Code Documentation

- Add godoc comments for all exported functions, types, and constants
- Include examples in documentation where helpful
- Keep comments concise but informative

```go
// GetUserBalance retrieves the BOGO token balance for a given address.
// It returns the balance as a big.Int and any error encountered.
//
// Example:
//   balance, err := GetUserBalance("0x1234...")
//   if err != nil {
//       return err
//   }
func GetUserBalance(address string) (*big.Int, error) {
    // Implementation
}
```

### API Documentation

- Update OpenAPI specification for new endpoints
- Include request/response examples
- Document error responses

## 🔐 Security Guidelines

### Sensitive Data

- Never commit private keys, secrets, or credentials
- Use environment variables for configuration
- Sanitize all user inputs
- Validate all blockchain addresses

### Security Testing

- Test for common vulnerabilities
- Validate input parameters
- Check authentication and authorization
- Test rate limiting

## 🚀 Deployment

### Testing Deployment

Before submitting changes that affect deployment:

1. **Test Docker build**
   ```bash
   docker build -t bogowi-api-test .
   docker run -p 3001:3001 --env-file .env bogowi-api-test
   ```

2. **Verify health endpoints**
   ```bash
   curl http://localhost:3001/api/health
   ```

## 📋 Pull Request Process

### Before Submitting

1. **Ensure your code follows standards**
   - All tests pass
   - Code is formatted and vetted
   - Documentation is updated

2. **Rebase on latest develop**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout your-feature-branch
   git rebase develop
   ```

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes or documented properly
```

### Review Process

1. **Automated Checks**: All CI checks must pass
2. **Code Review**: At least one reviewer approval required
3. **Testing**: Manual testing for significant changes
4. **Documentation**: Ensure documentation is updated

## 🐛 Bug Reports

### Before Reporting

1. Check existing issues
2. Verify the bug on latest version
3. Gather reproduction steps

### Bug Report Template

```markdown
**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Environment:**
- OS: [e.g., macOS, Linux]
- Go version: [e.g., 1.21.0]
- API version: [e.g., v1.0.0]

**Additional context**
Add any other context about the problem here.
```

## 🆕 Feature Requests

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear and concise description of what the problem is.

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.
```

## 📞 Communication

### Channels

- **Issues**: Bug reports and feature requests
- **Discussions**: General questions and discussions
- **Pull Requests**: Code contributions

### Guidelines

- Be respectful and constructive
- Provide clear and detailed information
- Search existing issues before creating new ones
- Use appropriate labels and templates

## 🎯 Development Priorities

### Current Focus Areas

1. **Performance Optimization**
   - Reduce API response times
   - Optimize database queries
   - Improve caching strategies

2. **Security Enhancements**
   - Enhanced authentication
   - Input validation improvements
   - Security audit recommendations

3. **Feature Development**
   - WebSocket support
   - Advanced analytics
   - Multi-chain compatibility

### Getting Started Recommendations

Good first issues for new contributors:

- Documentation improvements
- Unit test additions
- Code quality improvements
- Small bug fixes

## 📜 License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to BOGOWI Blockchain API! Your efforts help build a sustainable future through blockchain technology. 🌊🌍
