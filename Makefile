# Makefile for proxy-tunnel

# Build tags for all features
TAGS = with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api

# Default target
all: build

# Build with all features enabled
build:
	@echo "Building proxy-tunnel with all features..."
	go build -tags "$(TAGS)" -o proxy-tunnel ./cmd/proxy-tunnel/
	@echo "✓ Build complete: proxy-tunnel (with uTLS, QUIC, WireGuard, DHCP, Clash API)"

# Build minimal version (only uTLS for Reality support)
build-minimal:
	@echo "Building minimal proxy-tunnel..."
	go build -tags "with_utls" -o proxy-tunnel ./cmd/proxy-tunnel/
	@echo "✓ Build complete: proxy-tunnel (minimal with uTLS)"

# Build without any optional features (smallest binary)
build-basic:
	@echo "Building basic proxy-tunnel..."
	go build -o proxy-tunnel ./cmd/proxy-tunnel/
	@echo "✓ Build complete: proxy-tunnel (basic)"

# Get dependencies
deps:
	@echo "Getting dependencies..."
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f proxy-tunnel
	@echo "✓ Clean complete"

# Run Go tests
test:
	@echo "Running Go tests..."
	go test -v
	@echo "✓ All tests passed"

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -cover
	@echo "✓ Tests complete"

# Run integration test (binary)
test-binary: build-minimal
	@echo "Testing proxy startup..."
	@timeout 5 ./proxy-tunnel -link 'http://example.com:8080' > /tmp/proxy-test.log 2>&1 & \
	sleep 2; \
	if grep -q "sing-box started" /tmp/proxy-test.log; then \
		echo "✓ Binary test passed"; \
		pkill -f proxy-tunnel; \
	else \
		echo "✗ Binary test failed"; \
		cat /tmp/proxy-test.log; \
		exit 1; \
	fi

# Install to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin..."
	sudo cp proxy-tunnel /usr/local/bin/
	@echo "✓ Installed successfully"

# Uninstall from /usr/local/bin
uninstall:
	@echo "Uninstalling..."
	sudo rm -f /usr/local/bin/proxy-tunnel
	@echo "✓ Uninstalled successfully"

# Show build info
info:
	@echo "Build Configuration:"
	@echo "  Tags: $(TAGS)"
	@echo ""
	@echo "Available targets:"
	@echo "  make build         - Build with all features (default)"
	@echo "  make build-minimal - Build with uTLS only"
	@echo "  make build-basic   - Build basic version"
	@echo "  make deps          - Download dependencies"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make test          - Run Go tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make test-binary   - Test the compiled binary"
	@echo "  make install       - Install to /usr/local/bin"
	@echo "  make uninstall     - Remove from /usr/local/bin"

.PHONY: all build build-minimal build-basic deps clean test test-cover test-binary install uninstall info
