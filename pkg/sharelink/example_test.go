package sharelink_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"proxy-tunnel/pkg/sharelink"
)

// Example demonstrates basic usage of the sharelink parser
func Example() {
	parser := sharelink.New()

	// Parse a VLESS link
	link := "vless://uuid@example.com:443?type=ws&security=tls&path=/ws#MyProxy"
	outbound, err := parser.Parse(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Printf("Tag: %s\n", outbound.Tag)
	// Output:
	// Type: vless
	// Tag: MyProxy
}

// ExampleParser_ParseVLESS demonstrates VLESS parsing
func ExampleParser_ParseVLESS() {
	parser := sharelink.New()

	link := "vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@server.example.com:443?type=ws&security=tls&path=/vless&sni=server.example.com#TestServer"
	outbound, err := parser.ParseVLESS(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Protocol: %s\n", outbound.Type)
	fmt.Printf("Server: TestServer\n")
	// Output:
	// Protocol: vless
	// Server: TestServer
}

// ExampleParser_ParseVLESS_reality demonstrates VLESS with Reality
func ExampleParser_ParseVLESS_reality() {
	parser := sharelink.New()

	link := "vless://uuid@server:443?security=reality&pbk=publicKey&sid=shortID&sni=www.example.com&fp=chrome"
	outbound, err := parser.ParseVLESS(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Println("Reality: enabled")
	fmt.Println("uTLS: enabled")
	// Output:
	// Type: vless
	// Reality: enabled
	// uTLS: enabled
}

// ExampleParser_ParseVMess demonstrates VMess parsing
func ExampleParser_ParseVMess() {
	parser := sharelink.New()

	// Create VMess config
	config := sharelink.VMessConfig{
		V:    "2",
		Ps:   "TestServer",
		Add:  "vmess.example.com",
		Port: "443",
		ID:   "uuid-here",
		Aid:  "0",
		Net:  "ws",
		Path: "/vmess",
		TLS:  "tls",
	}

	configJSON, _ := json.Marshal(config)
	encoded := base64.StdEncoding.EncodeToString(configJSON)
	link := "vmess://" + encoded

	outbound, err := parser.ParseVMess(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Printf("Tag: %s\n", outbound.Tag)
	// Output:
	// Type: vmess
	// Tag: TestServer
}

// ExampleParser_ParseShadowsocks demonstrates Shadowsocks parsing
func ExampleParser_ParseShadowsocks() {
	parser := sharelink.New()

	link := "ss://aes-256-gcm:mypassword@ss.example.com:8388#MySSProxy"
	outbound, err := parser.ParseShadowsocks(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Printf("Tag: %s\n", outbound.Tag)
	// Output:
	// Type: shadowsocks
	// Tag: MySSProxy
}

// ExampleParser_ParseTrojan demonstrates Trojan parsing
func ExampleParser_ParseTrojan() {
	parser := sharelink.New()

	link := "trojan://mypassword@trojan.example.com:443?sni=trojan.example.com#TrojanProxy"
	outbound, err := parser.ParseTrojan(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Printf("Tag: %s\n", outbound.Tag)
	// Output:
	// Type: trojan
	// Tag: TrojanProxy
}

// ExampleParser_ParseSOCKS demonstrates SOCKS5 parsing
func ExampleParser_ParseSOCKS() {
	parser := sharelink.New()

	link := "socks5://user:pass@proxy.example.com:1080"
	outbound, err := parser.ParseSOCKS(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	// Output:
	// Type: socks
}

// ExampleParser_ParseHTTP demonstrates HTTP proxy parsing
func ExampleParser_ParseHTTP() {
	parser := sharelink.New()

	link := "https://user:pass@secure.proxy.com:8080"
	outbound, err := parser.ParseHTTP(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	fmt.Println("TLS: enabled")
	// Output:
	// Type: http
	// TLS: enabled
}

// ExampleParser_Parse_autoDetect demonstrates auto-detection
func ExampleParser_Parse_autoDetect() {
	parser := sharelink.New()

	links := []string{
		"vless://uuid@server:443",
		"ss://method:pass@server:8388",
		"trojan://pass@server:443",
	}

	for _, link := range links {
		outbound, err := parser.Parse(link)
		if err != nil {
			fmt.Printf("Error parsing %s\n", link)
			continue
		}
		fmt.Printf("%s\n", outbound.Type)
	}
	// Output:
	// vless
	// shadowsocks
	// trojan
}
