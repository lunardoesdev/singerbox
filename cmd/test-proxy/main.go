package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/lunardoesdev/singerbox"
	"github.com/sagernet/sing-box/option"
)

func main() {
	// Create proxy with direct outbound
	fmt.Println("Creating proxy...")
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound: option.Outbound{
			Type:    "direct",
			Tag:     "direct",
			Options: &option.DirectOutboundOptions{},
		},
		ListenAddr: "127.0.0.1:1080",
		LogLevel:   "info",
	})
	if err != nil {
		fmt.Printf("❌ Error creating proxy: %v\n", err)
		return
	}

	fmt.Println("Starting proxy...")
	err = pb.Start()
	if err != nil {
		fmt.Printf("❌ Error starting proxy: %v\n", err)
		return
	}
	defer pb.Stop()

	fmt.Println("✓ Proxy started successfully!")
	fmt.Printf("✓ Listening on: %s\n", pb.ListenAddr())

	// Wait a moment for it to fully initialize
	time.Sleep(500 * time.Millisecond)

	// Self-test the proxy
	fmt.Println("\n=== Running self-test ===")
	proxyURL, _ := url.Parse("http://127.0.0.1:1080")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		fmt.Printf("⚠️  HTTP request failed: %v\n", err)
		fmt.Println("   (This might be a network issue, not a proxy issue)")
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("✓ HTTP proxy working! Status: %s\n", resp.Status)
		fmt.Printf("  Response: %s\n", string(body))
	}

	fmt.Println("\n=== Proxy is ready for use ===")
	fmt.Println("Configure your browser or curl to use:")
	fmt.Println("  HTTP Proxy: 127.0.0.1:1080")
	fmt.Println("  SOCKS5 Proxy: 127.0.0.1:1080")
	fmt.Println("\nTest with: curl -x http://127.0.0.1:1080 http://example.com")
	fmt.Println("\nPress Ctrl+C to stop...")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nStopping proxy...")
}
