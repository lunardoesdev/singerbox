// Package sharelink provides functionality to parse proxy share links into sing-box outbound configurations.
package singerbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badoption"
)

// uuidRegex validates UUID format (with or without dashes)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{12}$`)

// MaxShareLinkLength is the maximum allowed length for share links (64KB)
const MaxShareLinkLength = 64 * 1024

// Parse parses a share link and returns a sing-box Outbound configuration
func Parse(link string) (option.Outbound, error) {
	link = strings.TrimSpace(link)

	if len(link) > MaxShareLinkLength {
		return option.Outbound{}, E.New("share link too long (max " + fmt.Sprint(MaxShareLinkLength) + " bytes)")
	}

	if strings.HasPrefix(link, "vless://") {
		return ParseVLESS(link)
	} else if strings.HasPrefix(link, "vmess://") {
		return ParseVMess(link)
	} else if strings.HasPrefix(link, "ss://") {
		return ParseShadowsocks(link)
	} else if strings.HasPrefix(link, "trojan://") {
		return ParseTrojan(link)
	} else if strings.HasPrefix(link, "socks://") || strings.HasPrefix(link, "socks5://") {
		return ParseSOCKS(link)
	} else if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return ParseHTTP(link)
	}

	return option.Outbound{}, E.New("unsupported protocol: " + strings.Split(link, "://")[0])
}

// ParseVLESS parses a VLESS share link
// Format: vless://uuid@server:port?params#name
func ParseVLESS(link string) (option.Outbound, error) {
	u, err := url.Parse(link)
	if err != nil {
		return option.Outbound{}, err
	}

	uuid := u.User.Username()
	if uuid == "" {
		return option.Outbound{}, E.New("missing UUID in VLESS link")
	}
	if !uuidRegex.MatchString(uuid) {
		return option.Outbound{}, E.New("invalid UUID format in VLESS link")
	}

	server := u.Hostname()
	if server == "" {
		return option.Outbound{}, E.New("missing server in VLESS link")
	}

	port := getPort(u.Host)
	query := u.Query()

	vlessOpts := option.VLESSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: uint16(port),
		},
		UUID: uuid,
	}

	// Parse transport type
	transport := query.Get("type")
	security := query.Get("security")

	// TLS configuration
	if security == "tls" || security == "reality" {
		tlsOpts := &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: query.Get("sni"),
		}

		if security == "reality" {
			pbk := query.Get("pbk")
			if pbk == "" {
				return option.Outbound{}, E.New("missing public key (pbk) for Reality in VLESS link")
			}

			// Reality requires uTLS
			tlsOpts.UTLS = &option.OutboundUTLSOptions{
				Enabled:     true,
				Fingerprint: "chrome", // Default fingerprint
			}
			tlsOpts.Reality = &option.OutboundRealityOptions{
				Enabled:   true,
				PublicKey: pbk,
				ShortID:   query.Get("sid"), // sid is optional
			}
			// Override fingerprint if specified
			if fp := query.Get("fp"); fp != "" {
				tlsOpts.UTLS.Fingerprint = fp
			}
		}

		vlessOpts.OutboundTLSOptionsContainer = option.OutboundTLSOptionsContainer{
			TLS: tlsOpts,
		}
	}

	// Transport configuration
	switch transport {
	case "ws":
		headers := make(badoption.HTTPHeader)
		if host := query.Get("host"); host != "" {
			headers["Host"] = []string{host}
		}
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    query.Get("path"),
				Headers: headers,
			},
		}
	case "grpc":
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "grpc",
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: query.Get("serviceName"),
			},
		}
	case "http", "h2":
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "http",
			HTTPOptions: option.V2RayHTTPOptions{
				Host: []string{query.Get("host")},
				Path: query.Get("path"),
			},
		}
	}

	if query.Get("flow") != "" {
		vlessOpts.Flow = query.Get("flow")
	}

	// Get tag from fragment (name)
	tag := "proxy"
	if u.Fragment != "" {
		tag = u.Fragment
	}

	return option.Outbound{
		Type:    "vless",
		Tag:     tag,
		Options: &vlessOpts,
	}, nil
}

// vmessConfig represents VMess JSON configuration (internal use only)
type vmessConfig struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
}

// ParseVMess parses a VMess share link
// Format: vmess://base64encoded
func ParseVMess(link string) (option.Outbound, error) {
	encoded := strings.TrimPrefix(link, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(encoded)
		if err != nil {
			return option.Outbound{}, E.New("invalid base64 encoding in VMess link")
		}
	}

	var vmess vmessConfig
	if err := json.Unmarshal(decoded, &vmess); err != nil {
		return option.Outbound{}, E.New("invalid JSON in VMess link: ", err)
	}

	if vmess.Add == "" {
		return option.Outbound{}, E.New("missing server address in VMess link")
	}
	if vmess.ID == "" {
		return option.Outbound{}, E.New("missing UUID in VMess link")
	}
	if !uuidRegex.MatchString(vmess.ID) {
		return option.Outbound{}, E.New("invalid UUID format in VMess link")
	}

	port := 443
	fmt.Sscanf(vmess.Port, "%d", &port)

	vmessOpts := option.VMessOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     vmess.Add,
			ServerPort: uint16(port),
		},
		UUID:     vmess.ID,
		Security: "auto",
	}

	// TLS
	if vmess.TLS == "tls" {
		vmessOpts.OutboundTLSOptionsContainer = option.OutboundTLSOptionsContainer{
			TLS: &option.OutboundTLSOptions{
				Enabled:    true,
				ServerName: vmess.SNI,
			},
		}
		if vmess.Host != "" && vmess.SNI == "" {
			vmessOpts.TLS.ServerName = vmess.Host
		}
	}

	// Transport
	switch vmess.Net {
	case "ws":
		headers := make(badoption.HTTPHeader)
		if vmess.Host != "" {
			headers["Host"] = []string{vmess.Host}
		}
		vmessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    vmess.Path,
				Headers: headers,
			},
		}
	case "grpc":
		vmessOpts.Transport = &option.V2RayTransportOptions{
			Type: "grpc",
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: vmess.Path,
			},
		}
	case "http", "h2":
		vmessOpts.Transport = &option.V2RayTransportOptions{
			Type: "http",
			HTTPOptions: option.V2RayHTTPOptions{
				Host: []string{vmess.Host},
				Path: vmess.Path,
			},
		}
	}

	tag := "proxy"
	if vmess.Ps != "" {
		tag = vmess.Ps
	}

	return option.Outbound{
		Type:    "vmess",
		Tag:     tag,
		Options: &vmessOpts,
	}, nil
}

// ParseShadowsocks parses a Shadowsocks share link
// Format: ss://base64encoded or ss://method:password@server:port
func ParseShadowsocks(link string) (option.Outbound, error) {
	return parseShadowsocksWithDepth(link, 0)
}

// parseShadowsocksWithDepth handles recursive decoding with depth limit
func parseShadowsocksWithDepth(link string, depth int) (option.Outbound, error) {
	const maxRecursionDepth = 2 // Allow one level of base64 decoding
	if depth > maxRecursionDepth {
		return option.Outbound{}, E.New("too many levels of base64 encoding in Shadowsocks link")
	}

	link = strings.TrimPrefix(link, "ss://")

	var method, password, server string
	var port int
	var tag string

	// Check for fragment (name)
	if idx := strings.Index(link, "#"); idx != -1 {
		tag = link[idx+1:]
		link = link[:idx]
	}

	if strings.Contains(link, "@") {
		parts := strings.Split(link, "@")
		userInfo := parts[0]
		serverInfo := parts[1]

		// Decode userinfo if base64
		decoded, err := base64.StdEncoding.DecodeString(userInfo)
		if err == nil {
			userInfo = string(decoded)
		} else {
			// Try URL decoding
			decoded, err = base64.URLEncoding.DecodeString(userInfo)
			if err == nil {
				userInfo = string(decoded)
			}
		}

		methodPass := strings.SplitN(userInfo, ":", 2)
		method = methodPass[0]
		if len(methodPass) > 1 {
			password = methodPass[1]
		}

		// Handle IPv6 server addresses
		if strings.HasPrefix(serverInfo, "[") {
			// IPv6: [::1]:port
			if idx := strings.LastIndex(serverInfo, "]:"); idx != -1 {
				server = serverInfo[1:idx]
				fmt.Sscanf(serverInfo[idx+2:], "%d", &port)
			} else if strings.HasSuffix(serverInfo, "]") {
				server = serverInfo[1 : len(serverInfo)-1]
			}
		} else {
			serverParts := strings.Split(serverInfo, ":")
			server = serverParts[0]
			if len(serverParts) > 1 {
				fmt.Sscanf(serverParts[1], "%d", &port)
			}
		}
	} else {
		// Entire string is base64 encoded
		decoded, err := base64.StdEncoding.DecodeString(link)
		if err != nil {
			decoded, err = base64.URLEncoding.DecodeString(link)
			if err != nil {
				return option.Outbound{}, E.New("invalid base64 encoding in Shadowsocks link")
			}
		}
		return parseShadowsocksWithDepth("ss://"+string(decoded), depth+1)
	}

	if server == "" {
		return option.Outbound{}, E.New("missing server in Shadowsocks link")
	}
	if method == "" {
		return option.Outbound{}, E.New("missing method in Shadowsocks link")
	}

	// Validate port
	if port < 1 || port > 65535 {
		port = 8388 // Default Shadowsocks port
	}

	if tag == "" {
		tag = "proxy"
	}

	return option.Outbound{
		Type: "shadowsocks",
		Tag:  tag,
		Options: &option.ShadowsocksOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     server,
				ServerPort: uint16(port),
			},
			Method:   method,
			Password: password,
		},
	}, nil
}

// ParseTrojan parses a Trojan share link
// Format: trojan://password@server:port?params#name
func ParseTrojan(link string) (option.Outbound, error) {
	u, err := url.Parse(link)
	if err != nil {
		return option.Outbound{}, err
	}

	password := u.User.Username()
	if password == "" {
		return option.Outbound{}, E.New("missing password in Trojan link")
	}

	server := u.Hostname()
	if server == "" {
		return option.Outbound{}, E.New("missing server in Trojan link")
	}

	port := getPort(u.Host)
	query := u.Query()

	trojanOpts := option.TrojanOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: uint16(port),
		},
		Password: password,
	}

	// TLS is usually enabled for Trojan
	sni := query.Get("sni")
	if sni == "" {
		sni = server
	}
	trojanOpts.OutboundTLSOptionsContainer = option.OutboundTLSOptionsContainer{
		TLS: &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: sni,
		},
	}

	// Transport
	transport := query.Get("type")
	switch transport {
	case "ws":
		headers := make(badoption.HTTPHeader)
		if host := query.Get("host"); host != "" {
			headers["Host"] = []string{host}
		}
		trojanOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path:    query.Get("path"),
				Headers: headers,
			},
		}
	case "grpc":
		trojanOpts.Transport = &option.V2RayTransportOptions{
			Type: "grpc",
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: query.Get("serviceName"),
			},
		}
	}

	tag := "proxy"
	if u.Fragment != "" {
		tag = u.Fragment
	}

	return option.Outbound{
		Type:    "trojan",
		Tag:     tag,
		Options: &trojanOpts,
	}, nil
}

// ParseSOCKS parses a SOCKS5 share link
// Format: socks5://[user:pass@]server:port
func ParseSOCKS(link string) (option.Outbound, error) {
	link = strings.TrimPrefix(link, "socks5://")
	link = strings.TrimPrefix(link, "socks://")

	u, err := url.Parse("socks5://" + link)
	if err != nil {
		return option.Outbound{}, err
	}

	server := u.Hostname()
	if server == "" {
		return option.Outbound{}, E.New("missing server in SOCKS link")
	}

	socksOpts := option.SOCKSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: uint16(getPortWithDefault(u.Host, 1080)),
		},
		Version: "5",
	}

	if u.User != nil {
		socksOpts.Username = u.User.Username()
		password, _ := u.User.Password()
		socksOpts.Password = password
	}

	return option.Outbound{
		Type:    "socks",
		Tag:     "proxy",
		Options: &socksOpts,
	}, nil
}

// ParseHTTP parses an HTTP/HTTPS proxy share link
// Format: http://[user:pass@]server:port or https://[user:pass@]server:port
func ParseHTTP(link string) (option.Outbound, error) {
	u, err := url.Parse(link)
	if err != nil {
		return option.Outbound{}, err
	}

	server := u.Hostname()
	if server == "" {
		return option.Outbound{}, E.New("missing server in HTTP link")
	}

	// Default port: 443 for HTTPS, 80 for HTTP
	defaultPort := 80
	if u.Scheme == "https" {
		defaultPort = 443
	}

	httpOpts := option.HTTPOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: uint16(getPortWithDefault(u.Host, defaultPort)),
		},
	}

	if u.User != nil {
		httpOpts.Username = u.User.Username()
		password, _ := u.User.Password()
		httpOpts.Password = password
	}

	if u.Scheme == "https" {
		httpOpts.OutboundTLSOptionsContainer = option.OutboundTLSOptionsContainer{
			TLS: &option.OutboundTLSOptions{
				Enabled: true,
			},
		}
	}

	return option.Outbound{
		Type:    "http",
		Tag:     "proxy",
		Options: &httpOpts,
	}, nil
}

// getPortWithDefault extracts port from host:port string, using specified default
// Handles IPv6 addresses like [::1]:8080
func getPortWithDefault(hostPort string, defaultPort int) int {
	// Handle IPv6 addresses
	if strings.HasPrefix(hostPort, "[") {
		// IPv6 format: [::1]:port or [::1]
		if idx := strings.LastIndex(hostPort, "]:"); idx != -1 {
			var port int
			if _, err := fmt.Sscanf(hostPort[idx+2:], "%d", &port); err == nil {
				if port >= 1 && port <= 65535 {
					return port
				}
			}
		}
		return defaultPort
	}

	// IPv4 or hostname format: host:port
	parts := strings.Split(hostPort, ":")
	if len(parts) < 2 {
		return defaultPort
	}
	var port int
	if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &port); err == nil {
		if port >= 1 && port <= 65535 {
			return port
		}
	}
	return defaultPort
}

// getPort extracts port from host:port string, returns 443 as default (for TLS protocols)
func getPort(hostPort string) int {
	return getPortWithDefault(hostPort, 443)
}
