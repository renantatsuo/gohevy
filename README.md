# gohevy

Go client for the Hevy app API.

This library is not officially affiliated with Hevy. Use at your own risk.
[Hevy API Documentation](https://api.hevyapp.com/docs).

## API key

To get your API key, you'll need a PRO subscription on the Hevy app. Then go to the [developer settings](https://hevy.com/settings?developer) and create a new API key.

## Installation

```bash
go get github.com/renantatsuo/gohevy
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/renantatsuo/gohevy"
)

func main() {
    // Create a new client with your API key
    client := hevy.NewClient("your-api-key-here")

    // Use context for request cancellation
    ctx := context.Background()

    // Get your workouts
    workouts, err := client.GetWorkouts(ctx, hevy.PaginationParams{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, workout := range workouts.Workouts {
        fmt.Printf("Workout: %s\n", workout.Title)
    }
}
```

## Configuration

### HTTP Client

```go
customHTTPClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

client := hevy.NewClient("your-api-key",
    hevy.WithHTTPClient(customHTTPClient),
)
```

## Error Handling

The library provides structured error handling with the `APIError` type:

```go
workouts, err := client.GetWorkouts(ctx, hevy.PaginationParams{Page: 1, PageSize: 10})
if err != nil {
    // Check if it's an API error
    if apiErr, ok := err.(*hevy.APIError); ok {
        fmt.Printf("Status Code: %d\n", apiErr.StatusCode)
        fmt.Printf("Status: %s\n", apiErr.Status)
        fmt.Printf("Response Body: %s\n", apiErr.Body)

        // Check error type
        if apiErr.IsClientError() {
            // 4xx error - bad request, auth failed, etc.
            fmt.Println("Client error")
        } else if apiErr.IsServerError() {
            // 5xx error - server issues, potentially retry-able
            fmt.Println("Server error")
        }
    } else {
        // Network error or other error
        fmt.Printf("Request error: %v\n", err)
    }
}
```

## License

MIT License.
