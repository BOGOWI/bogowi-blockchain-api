#!/bin/bash

# BOGOWI Blockchain API - Coverage Setup Script
# This script sets up gocov + gocov-html for enhanced coverage reporting

set -e

echo "🚀 Setting up Enhanced Coverage Tools for BOGOWI API..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Add GOPATH/bin to PATH if not already there
if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
    echo "📝 Adding Go bin directory to PATH..."
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Install gocov
echo "📦 Installing gocov..."
go install github.com/axw/gocov/gocov@latest

# Install gocov-html
echo "📦 Installing gocov-html..."
go install github.com/matm/gocov-html/cmd/gocov-html@latest

# Verify installations
echo "✅ Verifying installations..."
if command -v gocov &> /dev/null && command -v gocov-html &> /dev/null; then
    echo "✅ All tools installed successfully!"
    echo ""
    echo "🎯 Available commands:"
    echo "  make coverage-enhanced  - Generate enhanced HTML report"
    echo "  make coverage-all      - Generate all coverage reports"
    echo "  make coverage-install  - Reinstall coverage tools"
    echo ""
    echo "📊 Example usage:"
    echo "  make coverage-enhanced && open coverage-enhanced.html"
else
    echo "❌ Installation failed. Please check your Go installation."
    exit 1
fi

echo "🎉 Coverage setup complete!"
