package main

import (
	"fmt"
	"log"
	"net/http"

	"go/browser/lib"
)

func main() {
	startPort := 9223 // Different from our CDP proxy (9222)
	headless := true

	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("🌐 Chrome CDP Instance Launcher")
	fmt.Println("═══════════════════════════════════════════════════")

	instance, err := lib.LaunchChrome(startPort, headless)
	if err != nil {
		log.Fatalf("Failed to launch Chrome: %v", err)
	}
	defer instance.Close()

	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("✓ Chrome is running on port %d\n", startPort)
	if headless {
		fmt.Println("  Mode: Headless")
	} else {
		fmt.Println("  Mode: Headed (visible browser)")
	}
	fmt.Printf("📡 CDP Endpoint: http://localhost:%d\n", startPort)
	url, err := instance.GetWebSocketURL()
	if err != nil {
		fmt.Printf("Dev Tool: %s\n", url)
	}
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("Press Ctrl+C to stop Chrome...")
	fmt.Println()

	// Create a channel to receive messages
	webSocketUrlChannels := make(chan string, 1000)
	webSocketUrlChannels <- url

	// API Server - Register routes
	apiPort := ":8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lib.GetBrowserInstanceUrl(webSocketUrlChannels, w, r)
	})
	log.Printf("Server starting on http://localhost%s", apiPort)
	log.Println("\nAvailable endpoints:")
	log.Println("  GET  /                    - Hello World")

	if err := http.ListenAndServe(apiPort, nil); err != nil {
		log.Fatal(err)
	}

	// Keep running until interrupted
	select {}
}
