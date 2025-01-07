package main

import (
    "context"
    "fmt"
    "log"
    "moondream-go/internal/moondream"
    "os"
    "os/signal"
    "time"
)

func main() {
    if len(os.Args) < 3 {
        log.Fatal("Usage: moondream <function> <image-path> [options]\n" +
            "Functions: caption, query, detect, point\n" +
            "Example: moondream caption image.jpg")
    }

    apiKey := os.Getenv("mdAPI")
    if apiKey == "" {
        log.Fatal("Error: mdAPI environment variable not set")
    }

    // Create client with options
    client := moondream.NewMoondreamClient(
        apiKey,
        moondream.WithTimeout(30*time.Second),
    )

    // Create context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle interrupt signal
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    go func() {
        <-signalChan
        fmt.Println("\nReceived interrupt signal. Canceling operations...")
        cancel()
    }()

    function := os.Args[1]
    imagePath := os.Args[2]

    switch function {
    case "caption":
        caption, err := client.Caption(ctx, imagePath, "long")
        if err != nil {
            log.Fatalf("Error generating caption: %v", err)
        }
        fmt.Println("Caption:", caption)

    case "query":
        if len(os.Args) < 4 {
            log.Fatal("Usage: moondream query <image-path> <question>")
        }
        question := os.Args[3]
        answer, err := client.Query(ctx, imagePath, question)
        if err != nil {
            log.Fatalf("Error querying image: %v", err)
        }
        fmt.Println("Answer:", answer)

    case "detect":
        if len(os.Args) < 4 {
            log.Fatal("Usage: moondream detect <image-path> <object>")
        }
        object := os.Args[3]
        boundingBoxes, err := client.Detect(ctx, imagePath, object)
        if err != nil {
            log.Fatalf("Error detecting objects: %v", err)
        }
        fmt.Printf("Found %d instances of '%s':\n", len(boundingBoxes), object)
        for i, box := range boundingBoxes {
            fmt.Printf("  %d: x_min=%.2f, y_min=%.2f, x_max=%.2f, y_max=%.2f\n",
                i+1, box["x_min"], box["y_min"], box["x_max"], box["y_max"])
        }

    case "point":
        if len(os.Args) < 4 {
            log.Fatal("Usage: moondream point <image-path> <object>")
        }
        object := os.Args[3]
        points, err := client.Point(ctx, imagePath, object)
        if err != nil {
            log.Fatalf("Error pointing at objects: %v", err)
        }
        fmt.Printf("Found %d points for '%s':\n", len(points), object)
        for i, point := range points {
            fmt.Printf("  %d: x=%.2f, y=%.2f\n", i+1, point["x"], point["y"])
        }

    default:
        log.Fatalf("Unknown function: %s", function)
    }
}
