# 📊 Code Coverage for BOGOWI Blockchain API

This project uses **gocov + gocov-html** for enhanced test coverage reporting with beautiful, interactive HTML reports.

## 🚀 Quick Start

### Install Coverage Tools
```bash
# Option 1: Use our setup script
./scripts/setup-coverage.sh

# Option 2: Install manually
make coverage-install
```

### Generate Coverage Reports
```bash
# Generate enhanced HTML report (recommended)
make coverage-enhanced

# Generate all coverage reports (standard + enhanced)
make coverage-all

# Generate standard coverage report only
make test-coverage
```

## 📈 Available Reports

### Enhanced Report (Recommended) ⭐
- **File**: `coverage-enhanced.html`
- **Tool**: gocov + gocov-html
- **Features**: 
  - Beautiful, modern interface
  - Function-level coverage details
  - Package hierarchy visualization
  - Interactive navigation
  - Better styling and UX

### Standard Report
- **File**: `coverage.html` 
- **Tool**: Built-in `go tool cover`
- **Features**: Basic coverage visualization

## 🎯 Usage Examples

```bash
# Quick coverage check with auto-open
make coverage-enhanced && open coverage-enhanced.html

# Compare both reports
make coverage-all
open coverage.html coverage-enhanced.html

# CI/CD usage
make test-coverage  # Generates coverage.out for CI tools
```

## 📊 Current Coverage Status

- **Overall**: 14.4%
- **internal/config**: 42.9% ✅ Good
- **internal/sdk**: 18.3% 
- **internal/api**: 4.5%
- **main**: 0.0% (main functions)

## 🎨 Report Features

### Enhanced HTML Report Includes:
- **Package Overview**: High-level coverage metrics
- **File Browser**: Navigate through source files
- **Function Details**: Line-by-line coverage highlighting
- **Coverage Heatmap**: Visual representation of coverage density
- **Search & Filter**: Find specific functions or files
- **Responsive Design**: Works on desktop and mobile

### Standard Report Features:
- Basic line-by-line coverage
- Simple file navigation
- Coverage percentages

## 🔧 Configuration

Coverage tools are configured in:
- `Makefile`: Main coverage targets
- `scripts/setup-coverage.sh`: Automated setup
- `.gitignore`: Excludes coverage files from repository

## 📝 Tips

1. **Always use enhanced reports** for development - they're much more informative
2. **Set coverage goals** - aim for >70% on core business logic
3. **Focus on critical paths** - ensure high coverage on main workflows
4. **Exclude generated code** - use build tags or file patterns
5. **Integrate with CI** - use `coverage.out` for automated checks

## 🤝 Contributing

When adding new features:
1. Write tests first (TDD approach)
2. Run `make coverage-enhanced` to see coverage impact
3. Ensure critical paths have >80% coverage
4. Update tests if coverage drops significantly

## 📚 Resources

- [gocov GitHub](https://github.com/axw/gocov)
- [gocov-html GitHub](https://github.com/matm/gocov-html)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Coverage Story](https://blog.golang.org/cover)
