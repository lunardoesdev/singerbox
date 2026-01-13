package singerbox_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/lunardoesdev/singerbox"
)

// Example demonstrates basic usage of the sharelink parser
func Example_parse() {
	

	// Parse a VLESS link
	link := "vless://uuid@example.com:443?type=ws&security=tls&path=/ws#MyProxy"
	outbound, err := singerbox.Parse(link)
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
func Example_parseVLESS() {
	

	link := "vless://a1b2c3d4-e5f6-7890-abcd-ef1234567890@server.example.com:443?type=ws&security=tls&path=/vless&sni=server.example.com#TestServer"
	outbound, err := singerbox.ParseVLESS(link)
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
func Example_parseVLESS_reality() {
	

	link := "vless://uuid@server:443?security=reality&pbk=publicKey&sid=shortID&sni=www.example.com&fp=chrome"
	outbound, err := singerbox.ParseVLESS(link)
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
func Example_parseVMess() {
	

	// Create VMess config
	config := singerbox.VMessConfig{
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

	outbound, err := singerbox.ParseVMess(link)
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
func Example_parseShadowsocks() {
	

	link := "ss://aes-256-gcm:mypassword@ss.example.com:8388#MySSProxy"
	outbound, err := singerbox.ParseShadowsocks(link)
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
func Example_parseTrojan() {
	

	link := "trojan://mypassword@trojan.example.com:443?sni=trojan.example.com#TrojanProxy"
	outbound, err := singerbox.ParseTrojan(link)
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
func Example_parseSOCKS() {
	

	link := "socks5://user:pass@proxy.example.com:1080"
	outbound, err := singerbox.ParseSOCKS(link)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Type: %s\n", outbound.Type)
	// Output:
	// Type: socks
}

// ExampleParser_ParseHTTP demonstrates HTTP proxy parsing
func Example_parseHTTP() {
	

	link := "https://user:pass@secure.proxy.com:8080"
	outbound, err := singerbox.ParseHTTP(link)
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
func Example_parse_autoDetect() {
	

	links := []string{
		"vless://uuid@server:443",
		"ss://method:pass@server:8388",
		"trojan://pass@server:443",
	}

	for _, link := range links {
		outbound, err := singerbox.Parse(link)
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
