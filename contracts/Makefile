.PHONY: help test test-all coverage test-distributor test-treasury test-nft test-token clean compile

# Default target
help:
	@echo "BOGOWI Contracts - Test & Coverage Commands"
	@echo "==========================================="
	@echo ""
	@echo "Test Commands:"
	@echo "  make test-distributor    - Run BOGORewardDistributor tests"
	@echo "  make test-treasury       - Run MultisigTreasury tests"
	@echo "  make test-nft           - Run NFT contract tests"
	@echo "  make test-token         - Run BOGO token tests"
	@echo "  make test-all           - Run all tests"
	@echo "  make test               - Run all tests (alias)"
	@echo ""
	@echo "Coverage Commands:"
	@echo "  make coverage           - Run tests with coverage report"
	@echo "  make coverage-open      - Run coverage and open HTML report"
	@echo ""
	@echo "Other Commands:"
	@echo "  make compile            - Compile all contracts"
	@echo "  make clean              - Clean build artifacts and coverage"
	@echo "  make help               - Show this help menu"

# Individual test suites
test-distributor:
	@echo "Running BOGORewardDistributor tests..."
	@npx hardhat test test/BOGORewardDistributor.test.js

test-treasury:
	@echo "Running MultisigTreasury tests..."
	@npx hardhat test test/MultisigTreasury*.test.js

test-nft:
	@echo "Running NFT contract tests..."
	@npx hardhat test test/*NFT*.test.js

test-token:
	@echo "Running BOGO token tests..."
	@npx hardhat test test/BOGOToken*.test.js

# Run all tests
test-all:
	@echo "Running all tests..."
	@npx hardhat test

# Alias for test-all
test: test-all

# Coverage
coverage:
	@echo "Running tests with coverage..."
	@npx hardhat coverage

# Coverage with automatic browser opening
coverage-open: coverage
	@echo "Opening coverage report..."
	@open coverage/index.html || xdg-open coverage/index.html 2>/dev/null || echo "Please open coverage/index.html manually"

# Compile contracts
compile:
	@echo "Compiling contracts..."
	@npx hardhat compile

# Clean build artifacts and coverage
clean:
	@echo "Cleaning build artifacts and coverage..."
	@rm -rf artifacts cache coverage coverage.json
	@echo "Clean complete!"