package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Order Response struct
type Order struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// processOrder now accepts a String ID (Unique)
func processOrder(orderID string, baseURL string) string {
	url := fmt.Sprintf("%s/%s", baseURL, orderID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("❌ Network Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("⚠️ API Error (Status: %s)", resp.Status)
	}

	var o Order
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return "❌ JSON Decode Error"
	}

	return "✅"
}

// Worker now consumes string IDs
func worker(jobs <-chan string, results chan<- string, baseURL string, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		res := processOrder(id, baseURL)
		results <- res
	}
}

func main() {
	// LOAD TEST CONFIGURATION
	const totalOrders = 10000
	const concurrency = 100
	// Use 127.0.0.1 to avoid Windows IPv6 issues
	const apiURL = "http://127.0.0.1:8080/api/orders"

	// GENERATE A UNIQUE RUN ID (Timestamp based)
	// This ensures every test run creates unique keys in the DB
	runID := time.Now().Unix()

	fmt.Printf("🔥 STARTING FULL-STACK LOAD TEST\n")
	fmt.Printf("🆔 Test Run ID: %d\n", runID)
	fmt.Printf("🎯 Target API: %s\n", apiURL)
	fmt.Printf("📦 Total Orders: %d\n\n", totalOrders)

	// Change channel to String to support unique IDs
	jobs := make(chan string, totalOrders)
	results := make(chan string, totalOrders)
	var wg sync.WaitGroup

	start := time.Now()

	// 1. Start Workers
	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go worker(jobs, results, apiURL, &wg)
	}

	// 2. Dispatch Jobs with UNIQUE IDs
	go func() {
		for j := 1; j <= totalOrders; j++ {
			// Format: "RunID-OrderNumber" (e.g., "17012345-1")
			uniqueID := fmt.Sprintf("%d-%d", runID, j)
			jobs <- uniqueID
		}
		close(jobs)
	}()

	// 3. Wait for workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect Results
	successCount := 0
	failCount := 0

	for res := range results {
		if res == "✅" {
			successCount++
		} else {
			failCount++
			if failCount <= 5 {
				fmt.Printf("[Error Sample] %s\n", res)
			}
		}
	}

	elapsed := time.Since(start)
	rps := float64(totalOrders) / elapsed.Seconds()

	fmt.Printf("\n✨ LOAD TEST COMPLETED\n")
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("✅ Successful Orders:  %d\n", successCount)
	fmt.Printf("❌ Failed Requests:    %d\n", failCount)
	fmt.Printf("⏱️  Total Duration:     %v\n", elapsed)
	fmt.Printf("🚀 Throughput:         %.2f Requests/sec\n", rps)
	fmt.Printf("-----------------------------------\n")
}
