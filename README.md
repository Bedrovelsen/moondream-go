# Moondream Go Client

A Go client library and CLI tool for the [Moondream AI](https://moondream.ai) Moondream API, providing seamless access to Moondream's computer vision capabilities through a clean, idiomatic Go interface.

## Features

- **Visual Question Answering**: Ask natural language questions about images
- **Image Captioning**: Generate accurate and natural image captions
- **Object Detection**: Detect and locate objects in images
- **Object Pointing**: Get precise coordinate locations for objects in images

## Prerequisites

1. Get your API key from [Moondream Console](http://console.moondream.ai)
2. Set up your environment variable:
```bash
# Linux/macOS
export mdAPI=your-api-key

# Windows CMD
set mdAPI=your-api-key

# Windows PowerShell
$env:mdAPI="your-api-key"
```

## Installation

### Library

```bash
go get github.com/bedrovelsen/moondream-go
```

### CLI Tool

1. Clone the repository:
```bash
git clone https://github.com/bedrovelsen/moondream-go.git
cd moondream-go
```

2. Build for your platform:

```bash
# Build for your current platform
go build -o moondream ./cmd/moondream

# Or specify platform explicitly if needed:
# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o moondream ./cmd/moondream
# macOS (Intel/Apple Silicon)
GOOS=darwin go build -o moondream ./cmd/moondream
# Windows
GOOS=windows go build -o moondream.exe ./cmd/moondream
```

3. (Optional) Move the binary to your PATH:

```bash
# Linux/macOS
sudo mv moondream /usr/local/bin/

# Windows
# Move moondream.exe to a directory in your PATH
```

## Quick Start

### Using the Library

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/bedrovelsen/moondream-go/internal/moondream"
)

func main() {
    // Initialize client with your API key
    client := moondream.NewMoondreamClient(
        "your-api-key",
        moondream.WithTimeout(30*time.Second),
    )

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Example: Generate an image caption
    caption, err := client.Caption(ctx, "path/to/image.jpg", "normal", false)
    if err != nil {
        panic(err)
    }
    fmt.Println("Caption:", caption)
}
```

### Using the CLI

```bash
# Ask a question about an image
moondream query image.jpg "What is the main subject in this image?"

# Generate a natural image caption
moondream caption image.jpg

# Detect objects in an image
moondream detect image.jpg "person"

# Get precise coordinates of objects
moondream point image.jpg "coffee cup"
```

## API Reference

### Visual Question Answering (/query)
Ask natural language questions about images and receive detailed answers:
```go
answer, err := client.Query(ctx, imagePath, "What color is the car?")
```

### Image Captioning (/caption)
Generate accurate and natural image captions:
```go
// Default normal length caption
caption, err := client.Caption(ctx, imagePath, "normal", false)

// Short caption
caption, err := client.Caption(ctx, imagePath, "short", false)
```

### Object Detection (/detect)
Detect and locate objects in images with bounding boxes:
```go
boxes, err := client.Detect(ctx, imagePath, "person")
// Returns coordinates: x_min, y_min, x_max, y_max
```

### Object Pointing (/point)
Get precise coordinate locations for objects:
```go
points, err := client.Point(ctx, imagePath, "coffee cup")
// Returns coordinates: x, y
```

## Error Handling

The client provides detailed error information through custom error types:

```go
if err != nil {
    if apiErr, ok := err.(*moondream.APIError); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.StatusCode, apiErr.Message)
    }
    // Handle other error types
}
```

## Best Practices

1. **Always use context**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   ```

2. **Handle errors appropriately**:
   ```go
   if err != nil {
       log.Printf("Error: %v", err)
       // Handle specific error types
   }
   ```

3. **Clean up resources**:
   ```go
   defer cancel() // Cancel context when done
   ```

## TODO

Future improvments

1. **Streaming Support**
   - Add streaming response support for caption and query endpoints
   - Implement using Go channels for idiomatic streaming

2. **Input Flexibility**
   - Support `io.Reader` interface for image input
   - Support `image.Image` interface from the standard library
   - Add helper functions for common image formats

3. **Advanced Features**
   - Consider local inference support
   - Add proper response streaming

4. **Testing**
   - Add comprehensive unit tests
   - Add integration tests
   - Add example tests for documentation

5. **Documentation**
   - Add godoc examples
   - Document supported image formats and size limits
   - Add more usage examples

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.


## Acknowledgments

- Thanks to the Moondream team for providing the API
- Inspired by Go best practices and idioms
