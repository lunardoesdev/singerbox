# Quick Start Guide

## The Fixes

### Fix 1: Missing endpoint registry in context
The initial error **"missing endpoint registry in context"** was resolved by importing the `include` package which registers all protocol handlers with sing-box.

### Fix 2: Interface conversion panic
The second error **"interface conversion: interface {} is option.HTTPMixedInboundOptions, not *option.HTTPMixedInboundOptions"** was resolved by using pointers for all Options fields in Inbound and Outbound structs.

### Fix 3: uTLS not included
The third error **"uTLS, which is required by reality is not included in this build"** was resolved by building with the `with_utls` build tag (and other optional feature tags).

## What Changed

### Change 1: Initialize context with protocol registries

```go
import "github.com/sagernet/sing-box/include"

// In main():
ctx := context.Background()
ctx = include.Context(ctx) // Register all protocol handlers
```

### Change 2: Use pointers for Options

All inbound and outbound Options must be pointers:

```go
// Correct ✓
Options: &option.HTTPMixedInboundOptions{...}

// Wrong ✗
Options: option.HTTPMixedInboundOptions{...}
```

### Change 3: Build with feature tags

Build command now includes necessary tags:

```bash
# Recommended (all features)
go build -tags "with_quic,with_utls,with_wireguard,with_dhcp,with_clash_api" -o proxy-tunnel

# Or use Makefile
make build
```

## Usage

### 1. Basic Example (SOCKS5 Upstream)

If you have an existing SOCKS5 proxy (like Tor on port 9050):

```bash
./proxy-tunnel -link 'socks5://127.0.0.1:9050'
```

This creates:
- Local SOCKS5 proxy: `127.0.0.1:1080`
- Local HTTP proxy: `127.0.0.1:1081`

Both tunnel through your SOCKS5 proxy.

### 2. HTTP/HTTPS Proxy

```bash
./proxy-tunnel -link 'http://proxy.example.com:8080'
```

With authentication:
```bash
./proxy-tunnel -link 'http://user:pass@proxy.example.com:8080'
```

### 3. Shadowsocks

```bash
./proxy-tunnel -link 'ss://aes-256-gcm:password@server.com:8388'
```

Or with base64 encoded credentials:
```bash
./proxy-tunnel -link 'ss://base64encodedstring@server.com:8388'
```

### 4. VLESS with WebSocket + TLS

```bash
./proxy-tunnel -link 'vless://uuid@server.com:443?type=ws&security=tls&path=/ws&host=server.com'
```

### 5. VLESS with Reality

```bash
./proxy-tunnel -link 'vless://uuid@server.com:443?security=reality&pbk=publickey&sid=shortid&sni=www.example.com'

# With custom uTLS fingerprint
./proxy-tunnel -link 'vless://uuid@server.com:443?security=reality&pbk=publickey&sid=shortid&sni=www.example.com&fp=chrome'
```

**Note**: Reality requires:
- Valid base64-encoded public key (pbk parameter)
- Short ID (sid parameter)
- SNI (Server Name Indication)
- uTLS is automatically enabled with "chrome" fingerprint (or specify with fp parameter)

### 6. VMess

```bash
./proxy-tunnel -link 'vmess://base64encodedconfig'
```

### 7. Trojan

```bash
./proxy-tunnel -link 'trojan://password@server.com:443?sni=server.com'
```

With WebSocket transport:
```bash
./proxy-tunnel -link 'trojan://password@server.com:443?type=ws&path=/ws&security=tls&sni=server.com'
```

## Testing Your Connection

Once the proxy is running, test it:

### Using curl
```bash
# Test HTTP proxy
curl -x http://127.0.0.1:1081 https://ifconfig.me

# Test SOCKS5 proxy
curl -x socks5://127.0.0.1:1080 https://ifconfig.me

# Verbose output
curl -v -x http://127.0.0.1:1081 https://api.ipify.org
```

### Using wget
```bash
# HTTP proxy
http_proxy=http://127.0.0.1:1081 wget -O- https://ifconfig.me

# SOCKS5 via tsocks or proxychains
```

### Browser Configuration

**Firefox:**
1. Settings → Network Settings
2. Manual proxy configuration
3. HTTP Proxy: `127.0.0.1`, Port: `1081`
   OR
4. SOCKS Host: `127.0.0.1`, Port: `1080`, SOCKS v5

**Chrome/Chromium:**
```bash
chromium --proxy-server="http://127.0.0.1:1081"
# or
chromium --proxy-server="socks5://127.0.0.1:1080"
```

## Custom Listen Address

Change the local proxy ports:

```bash
./proxy-tunnel -link 'your-link' -listen '0.0.0.0:8080' -http-port 8081
```

This listens on:
- SOCKS5: `0.0.0.0:8080` (accessible from other machines)
- HTTP: `0.0.0.0:8081`

## Troubleshooting

### "Error parsing share link"
- Check that your share link is properly formatted
- Make sure it starts with a supported protocol (vless://, vmess://, ss://, etc.)
- Try enclosing the link in single quotes

### "Error creating sing-box instance"
This was the original issue - make sure you're using the updated version with `include.Context()`.

### "Error starting sing-box"
- Check if the ports 1080/1081 are already in use
- Try using different ports with `-listen` and `-http-port`
- Verify your upstream proxy server is accessible

### Connection issues
- Test if the upstream proxy works directly first
- Check firewall rules
- Verify DNS resolution is working

## Performance Tips

1. **Use appropriate buffer sizes**: The default settings should work for most cases
2. **Network selection**: sing-box automatically selects the best network path
3. **Multiple instances**: You can run multiple instances on different ports

## Security Notes

1. **Local access only**: By default, proxies listen on 127.0.0.1 (localhost only)
2. **No authentication**: The local proxies don't require authentication
3. **Shared network**: If using 0.0.0.0, add firewall rules to restrict access
4. **Credentials in CLI**: Be careful when using passwords in command line (visible in process list)

## Advanced: Share Link Format Details

### VLESS
```
vless://UUID@server:port?type=transport&security=tls/reality&path=/path&host=host&sni=sni
```

Parameters:
- `type`: ws, grpc, http (transport)
- `security`: none, tls, reality
- `path`: WebSocket/HTTP path
- `host`: Host header
- `sni`: TLS server name
- `pbk`: Reality public key
- `sid`: Reality short ID
- `flow`: xtls-rprx-vision

### VMess (JSON base64 encoded)
```json
{
  "v": "2",
  "ps": "remark",
  "add": "server",
  "port": "443",
  "id": "uuid",
  "aid": "0",
  "net": "ws",
  "type": "none",
  "host": "example.com",
  "path": "/path",
  "tls": "tls",
  "sni": "example.com"
}
```

Then base64 encode: `vmess://base64encode(json)`

### Shadowsocks
```
ss://method:password@server:port
# or
ss://base64(method:password)@server:port
```

Common methods: aes-256-gcm, chacha20-poly1305, aes-128-gcm

## Exit

Press `Ctrl+C` to stop the proxy gracefully.

## Support

For issues with:
- **sing-box library**: https://github.com/SagerNet/sing-box
- **This tool**: Check the main README.md

## Next Steps

- Check README.md for full documentation
- See example.sh for more examples
- Explore sing-box features: https://sing-box.sagernet.org/
