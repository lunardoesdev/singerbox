# Makefile for proxy-tunnel

# Build tags for all features
TAGS = with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api

# Default target
all: build

# Build with all features enabled
build:
	@echo "Building proxy-tunnel with all features..."
	go build -tags "$(TAGS)" -o proxy-tunnel
	@echo "✓ Build complete: proxy-tunnel (with uTLS, QUIC, WireGuard, DHCP, Clash API)"

# Build minimal version (only uTLS for Reality support)
build-minimal:
	@echo "Building minimal proxy-tunnel..."
	go build -tags "with_utls" -o proxy-tunnel
	@echo "✓ Build complete: proxy-tunnel (minimal with uTLS)"

# Build without any optional features (smallest binary)
build-basic:
	@echo "Building basic proxy-tunnel..."
	go build -o proxy-tunnel
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

# Run tests
test:
	@echo "Testing proxy startup..."
	@timeout 5 ./proxy-tunnel -link 'http://example.com:8080' > /tmp/proxy-test.log 2>&1 & \
	sleep 2; \
	if grep -q "sing-box started" /tmp/proxy-test.log; then \
		echo "✓ Test passed"; \
		pkill -f proxy-tunnel; \
	else \
		echo "✗ Test failed"; \
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
	@echo "  make test          - Test the build"
	@echo "  make install       - Install to /usr/local/bin"
	@echo "  make uninstall     - Remove from /usr/local/bin"

.PHONY: all build build-minimal build-basic deps clean test install uninstall info
