package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	serverUrl = "http://localhost:8080"
	delay     = 10 * time.Millisecond
)

func main() {
	// CLI flags
	numUsers := flag.Int("users", 5, "Number of concurrent users to simulate")
	reqPerUser := flag.Int("requests", 10, "Number of requests each user will send")
	flag.Parse()

	fmt.Printf("Simulating %d users, %d requests each, delay = %v\n", *numUsers, *reqPerUser, delay)

	var wg sync.WaitGroup

	for i := 1; i <= *numUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			client := &http.Client{}

			for j := 1; j <= *reqPerUser; j++ {
				req, err := http.NewRequest("GET", serverUrl, nil)
				if err != nil {
					fmt.Printf("[User %d] Error creating request: %v\n", userID, err)
					continue
				}

				req.Header.Set("X-User-ID", fmt.Sprintf("User-%d", userID))

				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("[User %d] Request error: %v\n", userID, err)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("[User %d] Error reading response body: %v\n", userID, err)
				}
				resp.Body.Close()

				if resp.StatusCode == http.StatusTooManyRequests {
					fmt.Printf("[User %d, request no %d] (%d): %s", userID, j, resp.StatusCode, string(body))
				}
				time.Sleep(delay)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All users finished.")
}
