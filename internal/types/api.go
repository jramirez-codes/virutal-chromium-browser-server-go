// types/api.go
package types

// Data models
type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	WsUrl   string `json:"wsUrl"`
}
