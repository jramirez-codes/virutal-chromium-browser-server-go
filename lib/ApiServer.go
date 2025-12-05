package lib

import (
	"encoding/json"
	"net/http"
)

// Data models
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	WsUrl   string `json:"wsUrl"`
}

// 3. GET all users
func GetBrowserInstanceUrl(ch chan string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(ch) == 0 {
		response := Response{
			Success: false,
			Message: "Failed to retrieve WebSocket URL",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Message: "Users retrieved successfully",
		WsUrl:   <-ch,
	}

	json.NewEncoder(w).Encode(response)
}
