package moondream

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Default configuration values
const (
	defaultTimeout    = 30 * time.Second
	defaultBaseURL   = "https://api.moondream.ai/v1"
	defaultMaxRetries = 3
	defaultRetryDelay = 1 * time.Second
)

type MoondreamClient struct {
	config  ClientConfig
	client  *http.Client
	apiKey  string
}

// NewMoondreamClient creates a new client with the given API key and options
func NewMoondreamClient(apiKey string, opts ...ClientOption) *MoondreamClient {
	config := ClientConfig{
		Timeout:    defaultTimeout,
		BaseURL:    defaultBaseURL,
		MaxRetries: defaultMaxRetries,
		RetryDelay: defaultRetryDelay,
	}

	client := &MoondreamClient{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
		apiKey: apiKey,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *MoondreamClient) {
		c.config.Timeout = timeout
		c.client.Timeout = timeout
	}
}

// WithBaseURL sets the base URL for API requests
func WithBaseURL(baseURL string) ClientOption {
	return func(c *MoondreamClient) {
		c.config.BaseURL = baseURL
	}
}

func (c *MoondreamClient) encodeImage(imagePath string) (string, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Add data URI prefix for base64 encoded image
	return fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(imageData)), nil
}

func (c *MoondreamClient) sendRequest(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Moondream-Auth", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.config.RetryDelay * time.Duration(attempt)):
			}
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to send request: %w", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		if resp.StatusCode >= 400 {
			var apiErr APIError
			if err := json.Unmarshal(body, &apiErr); err != nil {
				apiErr = APIError{
					StatusCode: resp.StatusCode,
					Message:    string(body),
				}
			}
			lastErr = &apiErr
			if resp.StatusCode < 500 { // Don't retry client errors
				return lastErr
			}
			continue
		}

		if err := json.Unmarshal(body, result); err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (c *MoondreamClient) Caption(ctx context.Context, imagePath string, length string, stream bool) (string, error) {
	encodedImage, err := c.encodeImage(imagePath)
	if err != nil {
		return "", err
	}

	req := CaptionRequest{
		Image:  encodedImage,
		Length: length,
		Stream: stream,
	}

	var resp CaptionResponse
	if err := c.sendRequest(ctx, "/caption", req, &resp); err != nil {
		return "", err
	}

	return resp.Caption, nil
}

func (c *MoondreamClient) Query(ctx context.Context, imagePath string, question string) (string, error) {
	encodedImage, err := c.encodeImage(imagePath)
	if err != nil {
		return "", err
	}

	req := QueryRequest{
		Image:    encodedImage,
		Question: question,
	}

	var resp QueryResponse
	if err := c.sendRequest(ctx, "/query", req, &resp); err != nil {
		return "", err
	}

	return resp.Answer, nil
}

func (c *MoondreamClient) Detect(ctx context.Context, imagePath string, object string) ([]map[string]float64, error) {
	encodedImage, err := c.encodeImage(imagePath)
	if err != nil {
		return nil, err
	}

	req := DetectRequest{
		Image:  encodedImage,
		Object: object,
	}

	var resp DetectResponse
	if err := c.sendRequest(ctx, "/detect", req, &resp); err != nil {
		return nil, err
	}

	return resp.BoundingBoxes, nil
}

func (c *MoondreamClient) Point(ctx context.Context, imagePath string, object string) ([]map[string]float64, error) {
	encodedImage, err := c.encodeImage(imagePath)
	if err != nil {
		return nil, err
	}

	req := PointRequest{
		Image:  encodedImage,
		Object: object,
	}

	var resp PointResponse
	if err := c.sendRequest(ctx, "/point", req, &resp); err != nil {
		return nil, err
	}

	return resp.Points, nil
}
