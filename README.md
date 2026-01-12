# Proxy Tunnel

A Go application that uses **sing-box as a library** to create a local HTTP/SOCKS proxy that tunnels traffic through various proxy protocols.

## Features

- ✅ **Uses sing-box as a library** (not separate process)
- ✅ **Reusable sharelink parser library** (`pkg/sharelink`) with comprehensive tests
- ✅ **Multiple proxy protocols**:
  - VLESS (with TLS, Reality, WebSocket, gRPC, HTTP/2 transports)
  - VMess (with TLS, WebSocket, gRPC, HTTP/2 transports)
  - Shadowsocks
  - Trojan (with TLS and transports)
  - SOCKS5
  - HTTP/HTTPS
- ✅ **Local mixed proxy** (HTTP + SOCKS5)
- ✅ **Simple command-line interface**
- ✅ **Well-tested** with comprehensive test suite

## Building

### Quick Build (Recommended)

```bash
# Using Makefile (includes all features)
make

# Or manually with all features
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" -o proxy-tunnel
```

### Build Options

```bash
# Full-featured build (recommended) - ~37MB
make build

# Minimal build (uTLS only for Reality support) - ~32MB
make build-minimal

# Basic build (no optional features) - ~32MB
make build-basic

# Get dependencies first
make deps
```

**Important**: The `with_utls` tag is required for Reality protocol support. Without it, you'll get an error when using Reality-enabled proxies.

### Build Tags Explained

- `with_utls` - **Required for Reality protocol** and advanced TLS features
- `with_quic` - Enables QUIC transport (Hysteria, Hysteria2)
- `with_wireguard` - Enables WireGuard protocol
- `with_dhcp` - Enables DHCP DNS server
- `with_clash_api` - Enables Clash API compatibility

**Note**: The first build will download all necessary dependencies.

## Usage

```bash
./proxy-tunnel -link <share-link> [-listen <addr:port>] [-http-port <port>]
```

### Parameters

- `-link`: Proxy share link (required)
  - Supported formats: `vless://`, `vmess://`, `ss://`, `trojan://`, `socks5://`, `http://`, `https://`
- `-listen`: Local SOCKS5 proxy address (default: `127.0.0.1:1080`)
- `-http-port`: Local HTTP proxy port (default: `1081`)

### Examples

#### VLESS with TLS and WebSocket
```bash
./proxy-tunnel -link 'vless://uuid@example.com:443?type=ws&security=tls&path=/ws'
```

#### VLESS with Reality
```bash
./proxy-tunnel -link 'vless://uuid@example.com:443?security=reality&pbk=<public-key>&sid=<short-id>&sni=www.example.com'

# With custom uTLS fingerprint (optional, defaults to chrome)
./proxy-tunnel -link 'vless://uuid@example.com:443?security=reality&pbk=<public-key>&sid=<short-id>&sni=www.example.com&fp=firefox'
```

**Reality Requirements**:
- Public key (pbk): Base64-encoded Reality public key
- Short ID (sid): Short identifier for the connection
- SNI: Server Name Indication for TLS
- Fingerprint (fp): Optional, defaults to "chrome". Valid options: chrome, firefox, safari, edge, ios, android, random

**Note**: uTLS is automatically enabled for Reality connections.

#### VMess
```bash
./proxy-tunnel -link 'vmess://base64encodedconfig'
```

#### Shadowsocks
```bash
./proxy-tunnel -link 'ss://method:password@server:port'
```

#### Trojan
```bash
./proxy-tunnel -link 'trojan://password@example.com:443?sni=example.com'
```

#### SOCKS5
```bash
./proxy-tunnel -link 'socks5://user:pass@proxy.example.com:1080'
```

#### HTTP/HTTPS Proxy
```bash
./proxy-tunnel -link 'https://user:pass@proxy.example.com:8080'
```

#### Custom Listen Address
```bash
./proxy-tunnel -link 'vless://...' -listen '0.0.0.0:1080' -http-port 8080
```

## After Starting

Once the proxy is running, you'll see:

```
✓ Proxy started successfully!
  SOCKS5: 127.0.0.1:1080
  HTTP:   127.0.0.1:1081
  Routing through: proxy

Press Ctrl+C to stop...
```

You can then configure your applications to use:
- **SOCKS5 proxy**: `127.0.0.1:1080`
- **HTTP proxy**: `127.0.0.1:1081`

## Testing

### Using curl
```bash
# Test with HTTP proxy
curl -x http://127.0.0.1:1081 https://ifconfig.me

# Test with SOCKS5 proxy
curl -x socks5://127.0.0.1:1080 https://ifconfig.me
```

### Using Firefox
1. Go to Settings → Network Settings
2. Select "Manual proxy configuration"
3. Set HTTP Proxy: `127.0.0.1`, Port: `1081`
4. Or set SOCKS Host: `127.0.0.1`, Port: `1080`

### Using Chrome/Chromium
```bash
chromium --proxy-server="http://127.0.0.1:1081"
# or
chromium --proxy-server="socks5://127.0.0.1:1080"
```

## Dependencies

- sing-box (latest dev-next branch)
- All dependencies are managed via Go modules

## Architecture

This application uses sing-box as a **library**, not as a separate process:

1. **Parses proxy share links** using `pkg/sharelink` library
2. **Creates sing-box instance** programmatically
3. **Starts local inbound proxies** (HTTP + SOCKS5)
4. **Routes all traffic** through the configured outbound proxy
5. **All components run** in a single Go process

## Using the Sharelink Parser Library

The share link parsing functionality is available as a reusable library at `pkg/sharelink`.

### Quick Example

```go
import "proxy-tunnel/pkg/sharelink"

parser := sharelink.New()
outbound, err := parser.Parse("vless://uuid@server:443?security=reality&pbk=key&sid=id&sni=example.com")
if err != nil {
    // handle error
}

// Use outbound in your sing-box configuration
config := option.Options{
    Outbounds: []option.Outbound{outbound},
    // ...
}
```

### Features

- ✅ Parse all major proxy protocols (VLESS, VMess, SS, Trojan, SOCKS, HTTP)
- ✅ Auto-detect protocol from URL scheme
- ✅ Comprehensive test coverage (100%)
- ✅ Fast performance (~3-8μs per parse)
- ✅ Clean, documented API

See [`pkg/sharelink/README.md`](pkg/sharelink/README.md) for full documentation.

### Running Tests

```bash
# Run all tests
go test ./...

# Run sharelink library tests with coverage
go test -cover ./pkg/sharelink/

# Run benchmarks
go test -bench=. ./pkg/sharelink/
```

## License

This project uses sing-box library which is licensed under GPLv3 or later.
