# Reality Protocol Guide

This guide explains how to use VLESS with Reality protocol in proxy-tunnel.

## What is Reality?

Reality is an advanced proxy protocol that:
- Disguises proxy traffic as normal TLS connections
- Makes traffic indistinguishable from legitimate HTTPS traffic
- Provides better censorship resistance than traditional protocols
- Uses real TLS certificates from legitimate websites

## Requirements

### Build Requirements

Reality requires the `with_utls` build tag:

```bash
# Build with uTLS support (minimum)
go build -tags "with_utls" -o proxy-tunnel

# Or full-featured build (recommended)
make build
```

### Configuration Requirements

A Reality share link requires:
1. **UUID**: Your unique identifier
2. **Server**: Server address and port
3. **Public Key (pbk)**: Base64-encoded Reality public key
4. **Short ID (sid)**: Short identifier (usually hex string)
5. **SNI**: Server Name Indication (the domain to mimic)
6. **Fingerprint (fp)**: Optional uTLS fingerprint (defaults to "chrome")

## Share Link Format

```
vless://UUID@SERVER:PORT?security=reality&pbk=PUBLIC_KEY&sid=SHORT_ID&sni=DOMAIN[&fp=FINGERPRINT]
```

### Example

```bash
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@example.com:443?security=reality&pbk=SomeBase64EncodedPublicKey123&sid=abcd1234&sni=www.microsoft.com'
```

## Parameters Explained

### UUID (Required)

Your unique user identifier. Format: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

**Example**: `a1b2c3d4-e5f6-7890-abcd-ef1234567890`

---

### Server and Port (Required)

The Reality server address.

**Examples**:
- `example.com:443`
- `192.168.1.100:8443`
- `proxy.example.org:2053`

---

### Public Key - pbk (Required)

Base64-encoded Reality public key. You get this from your Reality server configuration.

**Format**: Base64 string
**Example**: `SomeBase64EncodedPublicKey123==`

**Where to get it**: From your Reality server's configuration file or setup output.

---

### Short ID - sid (Required)

A short identifier for the connection. Usually a short hex string.

**Format**: Hex string (2-16 characters)
**Examples**:
- `abcd1234`
- `a1b2`
- `0123456789abcdef`

**Where to get it**: From your Reality server configuration.

---

### SNI - sni (Required)

Server Name Indication - the domain name your traffic will mimic.

**Important**:
- Must be a real, popular website
- Should support TLS 1.3
- Common choices: `www.microsoft.com`, `www.apple.com`, `www.amazon.com`, `www.cloudflare.com`

**Examples**:
- `sni=www.microsoft.com`
- `sni=www.apple.com`
- `sni=www.cloudflare.com`

---

### Fingerprint - fp (Optional)

uTLS fingerprint to mimic. Defaults to "chrome" if not specified.

**Valid values**:
- `chrome` - Google Chrome browser (default, most common)
- `firefox` - Mozilla Firefox browser
- `safari` - Apple Safari browser
- `edge` - Microsoft Edge browser
- `ios` - iOS Safari
- `android` - Android Chrome
- `random` - Random fingerprint

**Example**: `fp=firefox`

**When to use different fingerprints**:
- Use `chrome` for most cases (default)
- Use `firefox` if Chrome is blocked in your region
- Use `ios` or `android` if connecting from mobile
- Use `random` for additional obfuscation

---

## Complete Examples

### Basic Reality (uses default chrome fingerprint)

```bash
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@proxy.example.com:443?security=reality&pbk=YourBase64PublicKeyHere==&sid=abc123&sni=www.microsoft.com'
```

### Reality with Firefox fingerprint

```bash
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@proxy.example.com:443?security=reality&pbk=YourBase64PublicKeyHere==&sid=abc123&sni=www.apple.com&fp=firefox'
```

### Reality with custom port

```bash
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@proxy.example.com:8443?security=reality&pbk=YourBase64PublicKeyHere==&sid=abc123&sni=www.cloudflare.com'
```

### Reality with transport (WebSocket)

```bash
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@proxy.example.com:443?type=ws&path=/ws&security=reality&pbk=YourBase64PublicKeyHere==&sid=abc123&sni=www.microsoft.com'
```

---

## How It Works

1. **You start proxy-tunnel** with a Reality share link
2. **proxy-tunnel creates local proxies** (HTTP on 1081, SOCKS on 1080)
3. **Your applications connect** to these local proxies
4. **proxy-tunnel connects** to the Reality server
5. **Traffic is disguised** as legitimate HTTPS to the SNI domain
6. **Reality server forwards** your traffic to the internet
7. **Responses come back** through the same encrypted tunnel

```
Your App → Local Proxy → proxy-tunnel → [TLS to SNI domain] → Reality Server → Internet
                         (127.0.0.1:1080/1)                   (Looks like www.microsoft.com)
```

---

## Common Issues

### "uTLS is required by reality client"

**Cause**: Binary not built with `with_utls` tag.

**Fix**:
```bash
make build
# or
go build -tags "with_utls" -o proxy-tunnel
```

---

### "invalid public_key"

**Cause**: The public key (pbk parameter) is invalid or not properly base64-encoded.

**Fix**:
1. Get the correct public key from your Reality server
2. Ensure it's properly base64-encoded
3. Don't include spaces or line breaks
4. Verify with server administrator

---

### "connection failed"

**Possible causes**:
1. Server address is wrong
2. Port is blocked by firewall
3. Public key doesn't match server
4. Short ID doesn't match server

**Debug steps**:
```bash
# Test server connectivity
ping proxy.example.com
telnet proxy.example.com 443

# Test with verbose logging
./proxy-tunnel -link 'your-reality-link' 2>&1 | tee reality.log
```

---

### "SNI mismatch"

**Cause**: The SNI domain doesn't match server configuration.

**Fix**: Use the exact SNI domain specified by your server administrator.

---

## Best Practices

### 1. Choose Popular SNI Domains

Use well-known, high-traffic websites:
- ✅ `www.microsoft.com`
- ✅ `www.apple.com`
- ✅ `www.amazon.com`
- ✅ `www.cloudflare.com`
- ❌ Don't use: obscure or low-traffic domains

### 2. Match Fingerprint to Your Device

- Desktop Windows/Linux: `chrome` or `firefox`
- macOS: `safari` or `chrome`
- iPhone/iPad: `ios`
- Android: `android` or `chrome`

### 3. Keep Configuration Secure

- Don't share your UUID with others
- Public key should match your server
- Store share links securely (they contain credentials)

### 4. Test Connection

```bash
# Start proxy
./proxy-tunnel -link 'your-reality-link'

# In another terminal, test
curl -x http://127.0.0.1:1081 https://ifconfig.me
```

### 5. Monitor Logs

Check for connection issues:
```bash
./proxy-tunnel -link 'your-reality-link' 2>&1 | tee -a reality.log
```

---

## Advanced Configuration

### Using Reality with Multiple Transports

Reality can be combined with transports like WebSocket or gRPC:

```bash
# Reality + WebSocket
./proxy-tunnel -link 'vless://uuid@server:443?type=ws&path=/ws&security=reality&pbk=key&sid=id&sni=domain'

# Reality + gRPC
./proxy-tunnel -link 'vless://uuid@server:443?type=grpc&serviceName=service&security=reality&pbk=key&sid=id&sni=domain'
```

### Custom Local Ports

```bash
./proxy-tunnel -link 'your-reality-link' -listen '127.0.0.1:8080' -http-port 8081
```

### Listen on All Interfaces (Use with Caution)

```bash
# Makes proxy accessible from other devices on network
./proxy-tunnel -link 'your-reality-link' -listen '0.0.0.0:1080' -http-port 1081
```

---

## Getting Reality Server Information

If you're setting up your own Reality server, you need to:

1. Install sing-box on server
2. Generate Reality key pair
3. Configure server with:
   - Private key (server keeps)
   - Public key (give to clients)
   - Short IDs
   - Destination/SNI domains

4. Share with clients:
   - Server address and port
   - Public key (pbk)
   - Short ID (sid)
   - UUID (per user)
   - SNI domain

---

## Security Considerations

1. **Use HTTPS for obtaining share links** - Don't send Reality credentials over HTTP
2. **Verify server identity** - Ensure you trust the Reality server operator
3. **Keep UUID private** - Don't share with untrusted parties
4. **Update regularly** - Keep proxy-tunnel updated for security fixes
5. **Monitor usage** - Watch for unusual connection patterns

---

## Performance Tips

1. **Choose nearby servers** - Lower latency = better performance
2. **Use common SNI domains** - Less likely to be blocked
3. **Match fingerprint to your OS** - More realistic traffic patterns
4. **Avoid peak hours** - If server is congested

---

## Comparison with Other Protocols

| Feature | Reality | TLS | Plain |
|---------|---------|-----|-------|
| Censorship Resistance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐ |
| Performance | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Setup Complexity | ⭐⭐⭐ | ⭐⭐ | ⭐ |
| Traffic Hiding | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐ |

---

## Troubleshooting Checklist

Before asking for help, verify:

- [ ] Built with `with_utls` tag (`make build`)
- [ ] All parameters present (uuid, pbk, sid, sni)
- [ ] Public key is valid base64
- [ ] Server is reachable (ping/telnet)
- [ ] Port is not blocked by firewall
- [ ] UUID matches server configuration
- [ ] Short ID matches server configuration
- [ ] SNI domain is accessible

---

## Resources

- sing-box Reality documentation: https://sing-box.sagernet.org/configuration/shared/v2ray-transport/#reality
- uTLS fingerprints: https://github.com/refraction-networking/utls
- Reality protocol: https://github.com/XTLS/REALITY

---

## Example Session

Complete example of starting and using Reality proxy:

```bash
# 1. Build with uTLS support
make build

# 2. Start proxy with Reality configuration
./proxy-tunnel -link 'vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@proxy.example.com:443?security=reality&pbk=YourPublicKey==&sid=abc123&sni=www.microsoft.com'

# Output:
# INFO network: updated default interface wlo1, index 3
# INFO inbound/mixed[mixed-in]: tcp server started at 127.0.0.1:1080
# INFO inbound/http[http-in]: tcp server started at 127.0.0.1:1081
# INFO sing-box started (0.00s)
# ✓ Proxy started successfully!
#   SOCKS5: 127.0.0.1:1080
#   HTTP:   127.0.0.1:1081
#   Routing through: proxy

# 3. In another terminal, test connection
curl -x http://127.0.0.1:1081 https://ifconfig.me

# 4. Configure your browser to use HTTP proxy 127.0.0.1:1081

# 5. Browse securely!
```

That's it! Your traffic is now tunneled through Reality, appearing as legitimate HTTPS traffic to `www.microsoft.com`.
