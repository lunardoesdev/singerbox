# Build Guide

This document explains how to build proxy-tunnel with different feature sets.

## Quick Start

```bash
# Easiest way - build with all features
make

# Or manually
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" -o proxy-tunnel
```

## Build Requirements

- Go 1.21 or later
- Linux, macOS, or Windows
- Internet connection (for first build to download dependencies)

## Build Tags

sing-box uses build tags to enable optional features. Here's what each tag does:

### `with_utls` ⭐ **REQUIRED FOR REALITY**

**Status**: Essential for most modern proxies
**Size Impact**: +0MB (minimal)
**Enables**:
- Reality protocol support
- Advanced TLS fingerprinting
- uTLS-based transports

**When to use**: Always include this unless you only use basic protocols (SOCKS5, HTTP, plain Shadowsocks).

**Example protocols that need this**:
- VLESS with Reality
- Any proxy with advanced TLS features

### `with_quic`

**Status**: Recommended
**Size Impact**: +5MB
**Enables**:
- Hysteria protocol
- Hysteria2 protocol
- QUIC-based transports
- HTTP/3 support

**When to use**: If you use Hysteria, Hysteria2, or any QUIC-based protocols.

### `with_wireguard`

**Status**: Optional
**Size Impact**: +2MB
**Enables**:
- WireGuard protocol support
- WireGuard-based transports

**When to use**: If you use WireGuard VPN or WireGuard-based proxy protocols.

### `with_dhcp`

**Status**: Optional
**Size Impact**: +1MB
**Enables**:
- DHCP DNS server
- Dynamic DNS configuration

**When to use**: For advanced networking setups with DHCP-based DNS.

### `with_clash_api`

**Status**: Optional
**Size Impact**: +1MB
**Enables**:
- Clash-compatible API
- Web dashboard support
- API-based configuration

**When to use**: If you need Clash API compatibility or web dashboard features.

## Build Variants

### Full-Featured Build (Recommended)

**Size**: ~37MB
**Command**:
```bash
make build
# or
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" -o proxy-tunnel
```

**Includes**: All protocols and features
**Best for**: Production use, maximum compatibility

### Minimal Build (uTLS Only)

**Size**: ~32MB
**Command**:
```bash
make build-minimal
# or
go build -tags "with_utls" -o proxy-tunnel
```

**Includes**: All basic protocols + Reality support
**Best for**: When you need Reality but want smaller binary

### Basic Build (No Optional Features)

**Size**: ~32MB
**Command**:
```bash
make build-basic
# or
go build -o proxy-tunnel
```

**Includes**: Only core protocols (VLESS, VMess, Shadowsocks, Trojan, SOCKS, HTTP)
**Limitations**:
- ❌ No Reality support
- ❌ No Hysteria/Hysteria2
- ❌ No WireGuard

**Best for**: Simple setups without Reality or QUIC protocols

## Protocol Requirements Matrix

| Protocol | Basic Build | + with_utls | + with_quic |
|----------|-------------|-------------|-------------|
| HTTP/HTTPS | ✅ | ✅ | ✅ |
| SOCKS5 | ✅ | ✅ | ✅ |
| Shadowsocks | ✅ | ✅ | ✅ |
| VMess | ✅ | ✅ | ✅ |
| VLESS (basic) | ✅ | ✅ | ✅ |
| VLESS + Reality | ❌ | ✅ | ✅ |
| Trojan | ✅ | ✅ | ✅ |
| Hysteria | ❌ | ❌ | ✅ |
| Hysteria2 | ❌ | ❌ | ✅ |
| WireGuard | Needs with_wireguard |
| SSH | ✅ | ✅ | ✅ |
| Tor | ✅ | ✅ | ✅ |

## Common Build Scenarios

### Scenario 1: Using VLESS with Reality

```bash
# Minimum required
go build -tags "with_utls" -o proxy-tunnel

# Recommended (full features)
make build
```

### Scenario 2: Using Hysteria2

```bash
# Minimum required
go build -tags "with_quic" -o proxy-tunnel

# Recommended
go build -tags "with_quic,with_utls" -o proxy-tunnel
```

### Scenario 3: Basic SOCKS5/HTTP only

```bash
# Basic build is sufficient
make build-basic
```

### Scenario 4: WireGuard VPN

```bash
go build -tags "with_wireguard,with_utls" -o proxy-tunnel
```

## Troubleshooting

### Error: "uTLS, which is required by reality is not included"

**Solution**: Build with `with_utls` tag:
```bash
go build -tags "with_utls" -o proxy-tunnel
```

### Error: "QUIC is not included in this build"

**Solution**: Build with `with_quic` tag:
```bash
go build -tags "with_quic" -o proxy-tunnel
```

### Binary size is too large

Try building without optional features:
```bash
# Minimal with just Reality support
make build-minimal

# Or completely basic
make build-basic
```

You can also strip debug symbols:
```bash
go build -ldflags="-s -w" -tags "with_utls" -o proxy-tunnel
```

### Build takes too long

First build downloads dependencies. Subsequent builds are much faster.

Speed up builds:
```bash
# Use build cache
go build -tags "with_utls" -o proxy-tunnel

# Parallel compilation (already enabled by default)
GOMAXPROCS=8 go build -tags "with_utls" -o proxy-tunnel
```

## Cross-Compilation

Build for different platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -tags "with_utls" -o proxy-tunnel-linux-amd64

# Linux ARM64 (Raspberry Pi, etc.)
GOOS=linux GOARCH=arm64 go build -tags "with_utls" -o proxy-tunnel-linux-arm64

# macOS AMD64 (Intel Mac)
GOOS=darwin GOARCH=amd64 go build -tags "with_utls" -o proxy-tunnel-darwin-amd64

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -tags "with_utls" -o proxy-tunnel-darwin-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -tags "with_utls" -o proxy-tunnel.exe

# Android ARM64
GOOS=android GOARCH=arm64 go build -tags "with_utls" -o proxy-tunnel-android
```

## Makefile Targets

```bash
make                # Build with all features (default)
make build          # Same as above
make build-minimal  # Build with uTLS only
make build-basic    # Build basic version
make deps           # Download dependencies
make clean          # Remove build artifacts
make test           # Test the build
make install        # Install to /usr/local/bin
make uninstall      # Remove from /usr/local/bin
make info           # Show build configuration
```

## Recommended Configurations

### For Most Users
```bash
make build
```
Includes all features, maximum compatibility.

### For Minimal Size + Reality Support
```bash
make build-minimal
```
Smaller binary while supporting modern protocols.

### For Production Servers
```bash
go build -ldflags="-s -w" -tags "with_quic,with_utls,with_wireguard" -o proxy-tunnel
```
Optimized binary with stripped debug symbols.

## Environment Variables

```bash
# Enable CGO (needed for some features)
CGO_ENABLED=1 go build -tags "with_utls" -o proxy-tunnel

# Disable CGO (static binary, but some features unavailable)
CGO_ENABLED=0 go build -tags "with_utls" -o proxy-tunnel

# Set Go module proxy (for faster downloads in some regions)
GOPROXY=https://goproxy.io,direct go build -tags "with_utls" -o proxy-tunnel
```

## Verification

After building, verify the binary works:

```bash
# Check version
./proxy-tunnel -h

# Test startup
timeout 3 ./proxy-tunnel -link 'http://example.com:8080'

# Or use make
make test
```

## Performance Optimization

For production builds with maximum performance:

```bash
go build \
  -ldflags="-s -w" \
  -trimpath \
  -tags "with_quic,with_utls,with_wireguard" \
  -o proxy-tunnel
```

Flags explained:
- `-ldflags="-s -w"` - Strip debug info (smaller binary)
- `-trimpath` - Remove file system paths from binary
- `-tags` - Enable optional features

This produces an optimized, production-ready binary.
