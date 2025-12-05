package main

import (
	"fmt"
	"log"
	"os"

	"go/browser/lib"
)

func main() {
	port := 9223 // Different from our CDP proxy (9222)
	headless := true

	// Check for --no-headless flag
	for _, arg := range os.Args[1:] {
		if arg == "--no-headless" || arg == "--headed" {
			headless = false
		}
	}

	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("🌐 Chrome CDP Instance Launcher")
	fmt.Println("═══════════════════════════════════════════════════")

	instance, err := lib.LaunchChrome(port, headless)
	if err != nil {
		log.Fatalf("Failed to launch Chrome: %v", err)
	}
	defer instance.Close()

	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("✓ Chrome is running on port %d\n", port)
	if headless {
		fmt.Println("  Mode: Headless")
	} else {
		fmt.Println("  Mode: Headed (visible browser)")
	}
	fmt.Printf("📡 CDP Endpoint: http://localhost:%d\n", port)
	url, err := instance.GetWebSocketURL()
	if err != nil {
		fmt.Printf("Dev Tool: %s\n", url)
	}
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("Press Ctrl+C to stop Chrome...")
	fmt.Println()

	// Keep running until interrupted
	select {}
}
