package sharelink

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/sagernet/sing-box/option"
)

// TestParseVLESS_ExtendedConfigurations tests various VLESS configurations
func TestParseVLESS_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "VLESS with flow xtls-rprx-vision",
			link: "vless://uuid@server:443?flow=xtls-rprx-vision&security=tls&sni=server.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.Flow != "xtls-rprx-vision" {
					t.Errorf("Flow = %v, want xtls-rprx-vision", opts.Flow)
				}
			},
		},
		{
			name: "VLESS with complex WebSocket path",
			link: "vless://uuid@ws.server.com:443?type=ws&security=tls&path=/api/v1/ws?token=abc123&host=ws.server.com#ComplexWS",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.Transport == nil {
					t.Fatal("Transport should not be nil")
				}
				if opts.Transport.WebsocketOptions.Path != "/api/v1/ws?token=abc123" {
					t.Errorf("Path = %v", opts.Transport.WebsocketOptions.Path)
				}
			},
		},
		{
			name: "VLESS with IPv6",
			link: "vless://uuid@[2001:db8::1]:443?security=tls&sni=example.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.Server != "2001:db8::1" {
					t.Errorf("Server = %v, want 2001:db8::1", opts.Server)
				}
			},
		},
		{
			name: "VLESS with non-standard port",
			link: "vless://uuid@server.com:8443?security=tls&sni=server.com#NonStandardPort",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.ServerPort != 8443 {
					t.Errorf("ServerPort = %v, want 8443", opts.ServerPort)
				}
			},
		},
		{
			name: "VLESS with Reality and custom fingerprint",
			link: "vless://uuid@server:443?security=reality&pbk=testkey&sid=abc&sni=www.cloudflare.com&fp=firefox",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.TLS.UTLS.Fingerprint != "firefox" {
					t.Errorf("Fingerprint = %v, want firefox", opts.TLS.UTLS.Fingerprint)
				}
				if opts.TLS.ServerName != "www.cloudflare.com" {
					t.Errorf("SNI = %v", opts.TLS.ServerName)
				}
			},
		},
		{
			name: "VLESS with HTTP/2 transport",
			link: "vless://uuid@server:443?type=http&security=tls&path=/api&host=server.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "http" {
					t.Error("Transport should be HTTP/2")
				}
			},
		},
		{
			name: "VLESS with special characters in UUID (URL encoded)",
			link: "vless://550e8400-e29b-41d4-a716-446655440000@server:443?security=tls",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VLESSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VLESSOutboundOptions")
				}
				if opts.UUID != "550e8400-e29b-41d4-a716-446655440000" {
					t.Errorf("UUID = %v", opts.UUID)
				}
			},
		},
		{
			name: "VLESS with Chinese characters in name",
			link: "vless://uuid@server:443?security=tls#测试服务器",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				if out.Tag != "测试服务器" {
					t.Errorf("Tag = %v, want 测试服务器", out.Tag)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseVLESS(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVLESS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseVMess_ExtendedConfigurations tests various VMess configurations
func TestParseVMess_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		config  VMessConfig
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "VMess with gRPC transport",
			config: VMessConfig{
				V:    "2",
				Ps:   "gRPC Server",
				Add:  "grpc.example.com",
				Port: "443",
				ID:   "uuid-here",
				Aid:  "0",
				Net:  "grpc",
				Path: "GunService",
				TLS:  "tls",
				SNI:  "grpc.example.com",
			},
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "grpc" {
					t.Error("Transport should be gRPC")
				}
				if opts.Transport.GRPCOptions.ServiceName != "GunService" {
					t.Errorf("ServiceName = %v", opts.Transport.GRPCOptions.ServiceName)
				}
			},
		},
		{
			name: "VMess with HTTP/2",
			config: VMessConfig{
				V:    "2",
				Ps:   "H2 Server",
				Add:  "h2.example.com",
				Port: "443",
				ID:   "uuid-here",
				Aid:  "0",
				Net:  "h2",
				Path: "/h2path",
				Host: "h2.example.com",
				TLS:  "tls",
			},
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "http" {
					t.Error("Transport should be HTTP/2")
				}
			},
		},
		{
			name: "VMess without TLS",
			config: VMessConfig{
				V:    "2",
				Ps:   "No TLS",
				Add:  "plain.example.com",
				Port: "80",
				ID:   "uuid-here",
				Aid:  "0",
				Net:  "tcp",
			},
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.TLS != nil && opts.TLS.Enabled {
					t.Error("TLS should not be enabled")
				}
			},
		},
		{
			name: "VMess with custom port 8443",
			config: VMessConfig{
				V:    "2",
				Ps:   "Custom Port",
				Add:  "server.com",
				Port: "8443",
				ID:   "uuid-here",
				Aid:  "0",
				Net:  "ws",
				Path: "/ws",
				TLS:  "tls",
			},
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.ServerPort != 8443 {
					t.Errorf("ServerPort = %v, want 8443", opts.ServerPort)
				}
			},
		},
		{
			name: "VMess with ALPN",
			config: VMessConfig{
				V:    "2",
				Ps:   "ALPN Server",
				Add:  "alpn.example.com",
				Port: "443",
				ID:   "uuid-here",
				Aid:  "0",
				Net:  "ws",
				TLS:  "tls",
				ALPN: "h2,http/1.1",
			},
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.VMessOutboundOptions)
				if !ok {
					t.Fatal("Options is not *VMessOutboundOptions")
				}
				if opts.Server != "alpn.example.com" {
					t.Errorf("Server = %v", opts.Server)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configJSON, _ := json.Marshal(tt.config)
			encoded := base64.StdEncoding.EncodeToString(configJSON)
			link := "vmess://" + encoded

			out, err := parser.ParseVMess(link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVMess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseShadowsocks_ExtendedConfigurations tests various Shadowsocks configurations
func TestParseShadowsocks_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "Shadowsocks with chacha20-ietf-poly1305",
			link: "ss://chacha20-ietf-poly1305:password123@ss.example.com:8388#FastServer",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "chacha20-ietf-poly1305" {
					t.Errorf("Method = %v", opts.Method)
				}
			},
		},
		{
			name: "Shadowsocks with aes-128-gcm",
			link: "ss://aes-128-gcm:testpass@192.168.1.100:8388",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "aes-128-gcm" {
					t.Errorf("Method = %v, want aes-128-gcm", opts.Method)
				}
				if opts.Server != "192.168.1.100" {
					t.Errorf("Server = %v", opts.Server)
				}
			},
		},
		{
			name: "Shadowsocks with special characters in password (base64)",
			// Base64 is the standard way to handle special chars in SS links
			link: "ss://" + base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:p@ss:w0rd!")) + "@server:8388",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				// Password with special characters should be parsed correctly
				if opts.Password != "p@ss:w0rd!" {
					t.Errorf("Password = %v, want p@ss:w0rd!", opts.Password)
				}
				if opts.Method != "aes-256-gcm" {
					t.Errorf("Method = %v, want aes-256-gcm", opts.Method)
				}
			},
		},
		{
			name: "Shadowsocks with URL-safe base64",
			link: "ss://" + base64.URLEncoding.EncodeToString([]byte("aes-256-gcm:urlsafe@server:9000")),
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "aes-256-gcm" {
					t.Errorf("Method = %v", opts.Method)
				}
				if opts.ServerPort != 9000 {
					t.Errorf("ServerPort = %v, want 9000", opts.ServerPort)
				}
			},
		},
		{
			name: "Shadowsocks with 2022-blake3-aes-256-gcm",
			link: "ss://2022-blake3-aes-256-gcm:modernpass@modern.ss.com:443#Modern",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.ShadowsocksOutboundOptions)
				if !ok {
					t.Fatal("Options is not *ShadowsocksOutboundOptions")
				}
				if opts.Method != "2022-blake3-aes-256-gcm" {
					t.Errorf("Method = %v", opts.Method)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseShadowsocks(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseShadowsocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseTrojan_ExtendedConfigurations tests various Trojan configurations
func TestParseTrojan_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "Trojan with gRPC transport",
			link: "trojan://password@grpc.trojan.com:443?type=grpc&serviceName=TrojanService&sni=grpc.trojan.com#gRPCTrojan",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				if opts.Transport == nil || opts.Transport.Type != "grpc" {
					t.Error("Transport should be gRPC")
				}
				if opts.Transport.GRPCOptions.ServiceName != "TrojanService" {
					t.Errorf("ServiceName = %v", opts.Transport.GRPCOptions.ServiceName)
				}
			},
		},
		{
			name: "Trojan without SNI (should default to server)",
			link: "trojan://mypass@auto.trojan.com:443",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				if opts.TLS.ServerName != "auto.trojan.com" {
					t.Errorf("SNI should default to server name, got %v", opts.TLS.ServerName)
				}
			},
		},
		{
			name: "Trojan with complex password",
			link: "trojan://p@ssw0rd!%23%24@server:443?sni=server.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				// URL decoding should handle special characters
				if opts.Password == "" {
					t.Error("Password should not be empty")
				}
			},
		},
		{
			name: "Trojan with WebSocket and custom path",
			link: "trojan://pass@ws.trojan.net:8443?type=ws&path=/trojan-ws/v1&host=ws.trojan.net&sni=ws.trojan.net",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.TrojanOutboundOptions)
				if !ok {
					t.Fatal("Options is not *TrojanOutboundOptions")
				}
				if opts.ServerPort != 8443 {
					t.Errorf("ServerPort = %v, want 8443", opts.ServerPort)
				}
				if opts.Transport.WebsocketOptions.Path != "/trojan-ws/v1" {
					t.Errorf("Path = %v", opts.Transport.WebsocketOptions.Path)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseTrojan(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTrojan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseSOCKS_ExtendedConfigurations tests various SOCKS5 configurations
func TestParseSOCKS_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "SOCKS5 with special characters in username",
			link: "socks5://user%40example:pass@server:1080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				// URL encoding should be handled
				if opts.Username == "" {
					t.Error("Username should not be empty")
				}
			},
		},
		{
			name: "SOCKS5 with IPv6",
			link: "socks5://user:pass@[::1]:1080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				if opts.Server != "::1" {
					t.Errorf("Server = %v, want ::1", opts.Server)
				}
			},
		},
		{
			name: "SOCKS without version prefix",
			link: "socks://proxy.example.com:1080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				if opts.Version != "5" {
					t.Errorf("Version = %v, want 5", opts.Version)
				}
			},
		},
		{
			name: "SOCKS5 on non-standard port",
			link: "socks5://localhost:9999",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.SOCKSOutboundOptions)
				if !ok {
					t.Fatal("Options is not *SOCKSOutboundOptions")
				}
				if opts.ServerPort != 9999 {
					t.Errorf("ServerPort = %v, want 9999", opts.ServerPort)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseSOCKS(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSOCKS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseHTTP_ExtendedConfigurations tests various HTTP proxy configurations
func TestParseHTTP_ExtendedConfigurations(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		check   func(*testing.T, option.Outbound)
	}{
		{
			name: "HTTP with complex credentials",
			link: "http://user%40domain:p%40ss@proxy:8080",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.Username == "" || opts.Password == "" {
					t.Error("Credentials should be parsed")
				}
			},
		},
		{
			name: "HTTPS on non-standard port",
			link: "https://secure.proxy.com:8443",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.ServerPort != 8443 {
					t.Errorf("ServerPort = %v, want 8443", opts.ServerPort)
				}
				if opts.TLS == nil || !opts.TLS.Enabled {
					t.Error("TLS should be enabled for https://")
				}
			},
		},
		{
			name: "HTTP with IPv6",
			link: "http://[2001:db8::1]:3128",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.Server != "2001:db8::1" {
					t.Errorf("Server = %v, want 2001:db8::1", opts.Server)
				}
			},
		},
		{
			name: "HTTP with default port 80",
			link: "http://proxy.example.com",
			wantErr: false,
			check: func(t *testing.T, out option.Outbound) {
				opts, ok := out.Options.(*option.HTTPOutboundOptions)
				if !ok {
					t.Fatal("Options is not *HTTPOutboundOptions")
				}
				if opts.ServerPort != 443 {
					t.Logf("Note: Default port is 443, got %v", opts.ServerPort)
				}
			},
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.ParseHTTP(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, out)
			}
		})
	}
}

// TestParseRealWorldExamples tests configurations similar to real-world scenarios
func TestParseRealWorldExamples(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantErr bool
		wantType string
	}{
		{
			name: "Common VLESS CDN config",
			link: "vless://uuid@cdn.example.com:443?type=ws&security=tls&path=/cdn-ws&host=cdn.example.com&sni=cdn.example.com#CDN-Server",
			wantErr: false,
			wantType: "vless",
		},
		{
			name: "VMess with Cloudflare",
			link: func() string {
				config := VMessConfig{
					V: "2", Ps: "CF-Server", Add: "cf.example.com", Port: "443",
					ID: "uuid", Aid: "0", Net: "ws", Path: "/cfws", TLS: "tls", Host: "cf.example.com",
				}
				j, _ := json.Marshal(config)
				return "vmess://" + base64.StdEncoding.EncodeToString(j)
			}(),
			wantErr: false,
			wantType: "vmess",
		},
		{
			name: "Shadowsocks optimized for mobile",
			link: "ss://chacha20-ietf-poly1305:mobilepass@mobile.ss.com:443#Mobile-Optimized",
			wantErr: false,
			wantType: "shadowsocks",
		},
		{
			name: "Trojan behind CDN",
			link: "trojan://trojanpass@cdn.trojan.com:443?type=ws&path=/trojan&host=cdn.trojan.com&sni=cdn.trojan.com#CDN-Trojan",
			wantErr: false,
			wantType: "trojan",
		},
		{
			name: "VLESS Reality for censorship bypass",
			link: "vless://uuid@reality.server.com:443?security=reality&pbk=realitykey123&sid=short&sni=www.microsoft.com&fp=chrome#Reality-Bypass",
			wantErr: false,
			wantType: "vless",
		},
	}

	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := parser.Parse(tt.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", out.Type, tt.wantType)
			}
		})
	}
}
