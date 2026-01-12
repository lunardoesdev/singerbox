package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lunardoesdev/singerbox"
)

var (
	shareLink  = flag.String("link", "", "Share link (vless://, vmess://, ss://, etc.)")
	listenAddr = flag.String("listen", "127.0.0.1:1080", "Local proxy listen address")
	httpPort   = flag.Int("http-port", 1081, "HTTP proxy port")
)

func main() {
	flag.Parse()

	if *shareLink == "" {
		fmt.Println("Usage: proxy-tunnel -link <share-link> [-listen <addr:port>] [-http-port <port>]")
		fmt.Println("\nSupported protocols: vless, vmess, ss, trojan, http, socks")
		fmt.Println("\nExample:")
		fmt.Println("  proxy-tunnel -link 'vless://uuid@server:443?type=ws&security=tls'")
		os.Exit(1)
	}

	// Parse share link using the library
	parser := singerbox.NewParser()
	outbound, err := parser.Parse(*shareLink)
	if err != nil {
		fmt.Printf("Error parsing share link: %v\n", err)
		os.Exit(1)
	}

	// Create proxy box
	pb, err := singerbox.NewProxyBox(singerbox.ProxyBoxConfig{
		Outbound:   outbound,
		ListenAddr: *listenAddr,
		HTTPPort:   *httpPort,
		LogLevel:   "info",
	})
	if err != nil {
		fmt.Printf("Error creating proxy box: %v\n", err)
		os.Exit(1)
	}

	// Start the proxy
	err = pb.Start()
	if err != nil {
		fmt.Printf("Error starting proxy: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Proxy started successfully!\n")
	fmt.Printf("  SOCKS5: %s\n", pb.ListenAddr())
	fmt.Printf("  HTTP:   %s\n", pb.HTTPAddr())
	fmt.Printf("  Routing through: %s\n", outbound.Tag)
	fmt.Println("\nPress Ctrl+C to stop...")

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nStopping proxy...")
	pb.Stop()
}
