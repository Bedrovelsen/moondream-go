package moondream

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewMoondreamClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewMoondreamClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.config.BaseURL != defaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", defaultBaseURL, client.config.BaseURL)
	}
}

func TestClientOptions(t *testing.T) {
	apiKey := "test-api-key"
	customTimeout := 5 * time.Second
	customBaseURL := "https://custom.api.example.com"

	client := NewMoondreamClient(
		apiKey,
		WithTimeout(customTimeout),
		WithBaseURL(customBaseURL),
	)

	if client.config.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.config.Timeout)
	}

	if client.config.BaseURL != customBaseURL {
		t.Errorf("Expected base URL %s, got %s", customBaseURL, client.config.BaseURL)
	}
}

func TestEncodeImage(t *testing.T) {
	client := NewMoondreamClient("test-api-key")

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write some dummy image data
	testData := []byte("test image data")
	if _, err := tmpfile.Write(testData); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test encoding
	encoded, err := client.encodeImage(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Remove data URI prefix
	prefix := "data:image/jpeg;base64,"
	if !strings.HasPrefix(encoded, prefix) {
		t.Fatalf("Expected data URI prefix %q, got %q", prefix, encoded[:min(len(encoded), len(prefix))])
	}
	encoded = encoded[len(prefix):]

	// Decode and verify
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if string(decoded) != string(testData) {
		t.Errorf("Expected decoded data %s, got %s", string(testData), string(decoded))
	}
}

func TestCaption(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-Moondream-Auth") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"caption": "test caption"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewMoondreamClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte("test image data")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test caption
	ctx := context.Background()

	// Test normal length
	caption, err := client.Caption(ctx, tmpfile.Name(), "normal", false)
	if err != nil {
		t.Errorf("Caption failed: %v", err)
	}
	if caption == "" {
		t.Error("Expected non-empty caption")
	}

	// Test short length
	caption, err = client.Caption(ctx, tmpfile.Name(), "short", false)
	if err != nil {
		t.Errorf("Caption failed: %v", err)
	}
	if caption == "" {
		t.Error("Expected non-empty caption")
	}
}

func TestCaptionWithTimeout(t *testing.T) {
	// Create a test server that delays
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"caption": "test caption"}`))
	}))
	defer server.Close()

	// Create client with test server URL and short timeout
	client := NewMoondreamClient(
		"test-api-key",
		WithBaseURL(server.URL),
		WithTimeout(100*time.Millisecond),
	)

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte("test image data")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test caption with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err = client.Caption(ctx, tmpfile.Name(), "long", false)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestQuery(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-Moondream-Auth") != "test-api-key" {
			t.Errorf("Expected API key header")
		}
		if r.URL.Path != "/query" {
			t.Errorf("Expected /query endpoint, got %s", r.URL.Path)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"answer": "test answer"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewMoondreamClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte("test image data")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test query
	ctx := context.Background()
	answer, err := client.Query(ctx, tmpfile.Name(), "What's in this image?")
	if err != nil {
		t.Fatal(err)
	}

	if answer != "test answer" {
		t.Errorf("Expected answer 'test answer', got '%s'", answer)
	}
}

func TestDetect(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-Moondream-Auth") != "test-api-key" {
			t.Errorf("Expected API key header")
		}
		if r.URL.Path != "/detect" {
			t.Errorf("Expected /detect endpoint, got %s", r.URL.Path)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"objects": [{"x": 0.5, "y": 0.5, "width": 0.3, "height": 0.4}]}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewMoondreamClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte("test image data")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test detect
	ctx := context.Background()
	boxes, err := client.Detect(ctx, tmpfile.Name(), "cat")
	if err != nil {
		t.Fatal(err)
	}

	if len(boxes) != 1 {
		t.Errorf("Expected 1 bounding box, got %d", len(boxes))
	}

	box := boxes[0]
	expectedValues := map[string]float64{
		"x":      0.5,
		"y":      0.5,
		"width":  0.3,
		"height": 0.4,
	}

	for key, expected := range expectedValues {
		if got := box[key]; got != expected {
			t.Errorf("Expected %s = %f, got %f", key, expected, got)
		}
	}
}

func TestPoint(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-Moondream-Auth") != "test-api-key" {
			t.Errorf("Expected API key header")
		}
		if r.URL.Path != "/point" {
			t.Errorf("Expected /point endpoint, got %s", r.URL.Path)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"points": [{"x": 0.3, "y": 0.7}]}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewMoondreamClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Create a temporary test image
	tmpfile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte("test image data")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test point
	ctx := context.Background()
	points, err := client.Point(ctx, tmpfile.Name(), "cat")
	if err != nil {
		t.Fatal(err)
	}

	if len(points) != 1 {
		t.Errorf("Expected 1 point, got %d", len(points))
	}

	point := points[0]
	expectedValues := map[string]float64{
		"x": 0.3,
		"y": 0.7,
	}

	for key, expected := range expectedValues {
		if got := point[key]; got != expected {
			t.Errorf("Expected %s = %f, got %f", key, expected, got)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
