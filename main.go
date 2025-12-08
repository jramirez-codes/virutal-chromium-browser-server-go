package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"virtual-browser/internal/api"
	"virtual-browser/internal/browser"
	"virtual-browser/internal/types"
	"virtual-browser/internal/util"
)

// Global Variables
var (
	instanceCloseMap = make(map[string]func() error)
	mu               sync.RWMutex
)
var wsURLChannels = make(chan string, 500)
var serverStats = types.ServerStatsResponse{
	StartTime:                 0,
	CPUUsage:                  0.0,
	MemoryUsage:               0,
	LiveChromeInstanceCount:   0,
	ServedChromeInstanceCount: 0,
}

func CreateInstance() (*browser.ChromeInstance, error) {
	// Get Available Port
	startPort, err := util.GetPort()
	if err != nil {
		return nil, err
	}

	// Get Chrome Instance
	instance, err := browser.LaunchChrome(startPort)
	if err != nil {
		return nil, err
	}

	// Get WebSocket URL
	wsURL, err := instance.GetWebSocketURL()
	if err != nil {
		return nil, err
	}

	// Add WebSocket URL to map
	mu.Lock()
	wsURLChannels <- wsURL
	instanceCloseMap[wsURL] = instance.Close
	mu.Unlock()

	return instance, nil
}

func StartAPIServer() {
	// API Server - Register routes
	apiPort := ":8080"

	// Get WebSocket URL
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		go api.GetBrowserInstanceUrl(wsURLChannels, w, r)

		// Create New Instance N+1 (Preload)
		go func() {
			_, err := CreateInstance()
			if err != nil {
				log.Fatalf("Failed to create instance: %v", err)
			}
		}()
	})

	// Kill WebSocket URL
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		api.KillBrowserInstance(&mu, instanceCloseMap, w, r)
		serverStats.ServedChromeInstanceCount++
	})

	// Get Server Stats
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		serverStats.LiveChromeInstanceCount = len(wsURLChannels)
		api.GetServerStats(&serverStats, w, r)
	})

	log.Printf("Server starting on http://localhost%s", apiPort)
	log.Println("\nAvailable endpoints:")
	log.Println(" GET - Get WebSocket URL")

	if err := http.ListenAndServe(apiPort, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Record Start Time
	serverStats.StartTime = time.Now().Unix()

	// Create Inital Instance N+1
	go func() {
		_, err := CreateInstance()
		if err != nil {
			log.Fatalf("Failed to create instance: %v", err)
		}
	}()

	// Start API Server
	go StartAPIServer()

	// Keep running until interrupted
	select {}
}
