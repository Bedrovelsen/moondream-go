package moondream

import (
	"fmt"
	"time"
)

// APIResponse represents a generic API response.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.StatusCode, e.Message)
}

// ClientConfig represents the configuration for the Moondream client
type ClientConfig struct {
	Timeout    time.Duration
	BaseURL    string
	MaxRetries int
	RetryDelay time.Duration
}

// ClientOption represents a function that modifies the client configuration
type ClientOption func(*MoondreamClient)

// CaptionRequest represents a request to generate an image caption
type CaptionRequest struct {
	Image  string `json:"image_url"`
	Length string `json:"length"`
}

// CaptionResponse represents the response from a caption request
type CaptionResponse struct {
	Caption string `json:"caption"`
}

// QueryRequest represents a request to query about an image
type QueryRequest struct {
	Image    string `json:"image_url"`
	Question string `json:"question"`
}

// QueryResponse represents the response from a query request
type QueryResponse struct {
	Answer string `json:"answer"`
}

// DetectRequest represents a request to detect objects in an image
type DetectRequest struct {
	Image  string `json:"image_url"`
	Object string `json:"object"`
}

// DetectResponse represents the response from a detect request
type DetectResponse struct {
	BoundingBoxes []map[string]float64 `json:"objects"`
}

// PointRequest represents a request to point at objects in an image
type PointRequest struct {
	Image  string `json:"image_url"`
	Object string `json:"object"`
}

// PointResponse represents the response from a point request
type PointResponse struct {
	Points []map[string]float64 `json:"points"`
}
