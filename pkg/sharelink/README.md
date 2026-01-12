# Sharelink Parser Library

A Go library for parsing proxy share links (vless, vmess, shadowsocks, trojan, socks5, http) into sing-box outbound configurations.

## Features

- ✅ **VLESS** - Supports TLS, Reality, WebSocket, gRPC, HTTP/2 transports
- ✅ **VMess** - Parses base64-encoded JSON configurations
- ✅ **Shadowsocks** - Supports both plain and base64-encoded formats
- ✅ **Trojan** - With TLS and transport support
- ✅ **SOCKS5** - With optional authentication
- ✅ **HTTP/HTTPS** - Basic and authenticated proxies
- ✅ **Auto-detection** - Automatically detects protocol from URL scheme
- ✅ **Comprehensive tests** - 100% test coverage
- ✅ **Fast** - Optimized for performance (~3-8μs per parse)

## Installation

```bash
go get proxy-tunnel/pkg/sharelink
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "proxy-tunnel/pkg/sharelink"
)

func main() {
    parser := sharelink.New()

    // Parse a VLESS link
    outbound, err := parser.Parse("vless://uuid@server:443?type=ws&security=tls")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Type: %s, Tag: %s\n", outbound.Type, outbound.Tag)
}
```

### Protocol-Specific Parsing

```go
parser := sharelink.New()

// VLESS
vlessOut, _ := parser.ParseVLESS("vless://uuid@server:443?security=reality&pbk=key&sid=id&sni=example.com")

// VMess
vmessOut, _ := parser.ParseVMess("vmess://base64encodedconfig")

// Shadowsocks
ssOut, _ := parser.ParseShadowsocks("ss://aes-256-gcm:password@server:8388")

// Trojan
trojanOut, _ := parser.ParseTrojan("trojan://password@server:443")

// SOCKS5
socksOut, _ := parser.ParseSOCKS("socks5://user:pass@server:1080")

// HTTP/HTTPS
httpOut, _ := parser.ParseHTTP("https://user:pass@server:8080")
```

### Auto-Detection

The `Parse()` method automatically detects the protocol:

```go
parser := sharelink.New()

links := []string{
    "vless://uuid@server:443",
    "vmess://base64...",
    "ss://method:pass@server:8388",
    "trojan://pass@server:443",
    "socks5://server:1080",
    "http://server:8080",
}

for _, link := range links {
    outbound, err := parser.Parse(link)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        continue
    }
    fmt.Printf("Parsed: %s\n", outbound.Type)
}
```

## Link Format Specifications

### VLESS

```
vless://UUID@SERVER:PORT?params#name

Parameters:
  - type: ws, grpc, http (transport)
  - security: tls, reality
  - path: WebSocket/HTTP path
  - host: Host header
  - sni: TLS server name
  - pbk: Reality public key
  - sid: Reality short ID
  - fp: uTLS fingerprint (chrome, firefox, safari, etc.)
  - flow: xtls-rprx-vision
```

**Example:**
```
vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?type=ws&security=tls&path=/ws&sni=example.com#MyServer
```

### VMess

```
vmess://BASE64_ENCODED_JSON

JSON structure:
{
  "v": "2",
  "ps": "remark",
  "add": "server",
  "port": "443",
  "id": "uuid",
  "aid": "0",
  "net": "ws|grpc|http|h2",
  "type": "none",
  "host": "example.com",
  "path": "/path",
  "tls": "tls|",
  "sni": "example.com"
}
```

### Shadowsocks

```
ss://METHOD:PASSWORD@SERVER:PORT#name
or
ss://BASE64(METHOD:PASSWORD)@SERVER:PORT#name

Common methods:
  - aes-256-gcm
  - chacha20-poly1305
  - aes-128-gcm
```

**Example:**
```
ss://aes-256-gcm:mypassword@server.com:8388#MyProxy
```

### Trojan

```
trojan://PASSWORD@SERVER:PORT?params#name

Parameters:
  - type: ws, grpc (optional transport)
  - path: transport path
  - host: Host header
  - sni: TLS server name
```

**Example:**
```
trojan://mypassword@example.com:443?sni=example.com#MyTrojan
```

### SOCKS5

```
socks5://[USER:PASS@]SERVER:PORT
```

**Examples:**
```
socks5://proxy.example.com:1080
socks5://user:pass@proxy.example.com:1080
```

### HTTP/HTTPS

```
http://[USER:PASS@]SERVER:PORT
https://[USER:PASS@]SERVER:PORT
```

**Examples:**
```
http://proxy.example.com:8080
https://user:pass@secure.proxy.com:8080
```

## API Reference

### type Parser

```go
type Parser struct{}
```

#### func New

```go
func New() *Parser
```

Creates a new Parser instance.

#### func (*Parser) Parse

```go
func (p *Parser) Parse(link string) (option.Outbound, error)
```

Parses a share link and auto-detects the protocol. Returns a sing-box Outbound configuration.

#### func (*Parser) ParseVLESS

```go
func (p *Parser) ParseVLESS(link string) (option.Outbound, error)
```

Parses a VLESS share link. Returns error if the link is invalid or missing required fields.

#### func (*Parser) ParseVMess

```go
func (p *Parser) ParseVMess(link string) (option.Outbound, error)
```

Parses a VMess share link (base64-encoded JSON).

#### func (*Parser) ParseShadowsocks

```go
func (p *Parser) ParseShadowsocks(link string) (option.Outbound, error)
```

Parses a Shadowsocks share link.

#### func (*Parser) ParseTrojan

```go
func (p *Parser) ParseTrojan(link string) (option.Outbound, error)
```

Parses a Trojan share link.

#### func (*Parser) ParseSOCKS

```go
func (p *Parser) ParseSOCKS(link string) (option.Outbound, error)
```

Parses a SOCKS5 share link.

#### func (*Parser) ParseHTTP

```go
func (p *Parser) ParseHTTP(link string) (option.Outbound, error)
```

Parses an HTTP/HTTPS proxy link.

## Error Handling

All parse methods return descriptive errors:

```go
parser := sharelink.New()
outbound, err := parser.Parse(link)
if err != nil {
    // Handle specific errors
    fmt.Printf("Parse error: %v\n", err)
    return
}
```

Common errors:
- `"unsupported protocol: xxx"` - Unknown protocol
- `"missing UUID in VLESS link"` - Required field missing
- `"invalid base64 encoding"` - Malformed base64
- `"invalid JSON in VMess link"` - Invalid VMess JSON

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./pkg/sharelink/

# Run tests with coverage
go test -cover ./pkg/sharelink/

# Run tests verbosely
go test -v ./pkg/sharelink/

# Run benchmarks
go test -bench=. -benchmem ./pkg/sharelink/
```

### Test Coverage

The library includes comprehensive tests for:
- ✅ All protocols (VLESS, VMess, SS, Trojan, SOCKS, HTTP)
- ✅ Various transports (WebSocket, gRPC, HTTP/2)
- ✅ Security options (TLS, Reality, uTLS)
- ✅ Error cases (invalid links, missing fields)
- ✅ Edge cases (base64 variants, optional parameters)
- ✅ Auto-detection logic
- ✅ Performance benchmarks

### Benchmark Results

```
BenchmarkParseVLESS-8          313958    3593 ns/op    1552 B/op    16 allocs/op
BenchmarkParseVMess-8          141874    8317 ns/op    1896 B/op    24 allocs/op
BenchmarkParseShadowsocks-8    616960    1719 ns/op     528 B/op    10 allocs/op
```

## Examples

### Example 1: Parse and Use with sing-box

```go
package main

import (
    "context"
    "proxy-tunnel/pkg/sharelink"
    "github.com/sagernet/sing-box"
    "github.com/sagernet/sing-box/option"
)

func main() {
    // Parse share link
    parser := sharelink.New()
    outbound, err := parser.Parse("vless://uuid@server:443?security=reality&pbk=key&sid=id&sni=example.com")
    if err != nil {
        panic(err)
    }

    // Create sing-box configuration
    config := option.Options{
        Outbounds: []option.Outbound{outbound},
        // ... rest of config
    }

    // Create sing-box instance
    instance, _ := box.New(box.Options{
        Context: context.Background(),
        Options: config,
    })

    _ = instance.Start()
}
```

### Example 2: Batch Parse Multiple Links

```go
func parseMultiple(links []string) ([]option.Outbound, []error) {
    parser := sharelink.New()
    var outbounds []option.Outbound
    var errors []error

    for _, link := range links {
        out, err := parser.Parse(link)
        if err != nil {
            errors = append(errors, err)
            continue
        }
        outbounds = append(outbounds, out)
    }

    return outbounds, errors
}
```

### Example 3: Validate Link Before Parsing

```go
func isValidLink(link string) bool {
    parser := sharelink.New()
    _, err := parser.Parse(link)
    return err == nil
}
```

## Dependencies

- `github.com/sagernet/sing-box` - For option types
- `github.com/sagernet/sing` - For common utilities

## License

This library is part of the proxy-tunnel project.

## Contributing

When contributing, please:
1. Add tests for new features
2. Ensure all tests pass
3. Update documentation
4. Follow Go best practices

## Support

For issues or questions:
- Check the test files for usage examples
- Review the protocol specifications above
- Open an issue in the main project repository
