package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Order Response struct (Must match the API response)
type Order struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}

// HTTP Client Configuration optimized for High Concurrency
// We use a custom Transport to pool connections properly
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100, // Keep up to 100 idle connections open
		MaxIdleConnsPerHost: 100, // Allow 100 connections to the specific host
		IdleConnTimeout:     90 * time.Second,
	},
}

// processOrder executes the HTTP request against the local API
func processOrder(orderID int, baseURL string) string {
	url := fmt.Sprintf("%s/%d", baseURL, orderID)

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

// Worker function: Consumes jobs from the channel and processes them
func worker(jobs <-chan int, results chan<- string, baseURL string, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jobs {
		res := processOrder(id, baseURL)
		results <- res
	}
}

func main() {
	// LOAD TEST CONFIGURATION
	// NOTE: 50,000 requests against a local DB is intense.
	// If you see errors, try lowering workers to 50 initially.
	const totalOrders = 10000
	const concurrency = 100
	const apiURL = "http://localhost:8080/api/orders"

	fmt.Printf("🔥 STARTING FULL-STACK LOAD TEST\n")
	fmt.Printf("🎯 Target API: %s\n", apiURL)
	fmt.Printf("📦 Total Orders: %d\n", totalOrders)
	fmt.Printf("⚡ Concurrent Workers: %d\n\n", concurrency)

	jobs := make(chan int, totalOrders)
	results := make(chan string, totalOrders)
	var wg sync.WaitGroup

	start := time.Now()

	// 1. Start Workers
	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go worker(jobs, results, apiURL, &wg)
	}

	// 2. Dispatch Jobs
	go func() {
		for j := 1; j <= totalOrders; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// 3. Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect Results
	successCount := 0
	failCount := 0

	// We iterate over the results channel as data comes in
	for res := range results {
		if res == "✅" {
			successCount++
		} else {
			failCount++
			// Print first 5 errors to help debugging, then silence
			if failCount <= 5 {
				fmt.Printf("[Error Sample] %s\n", res)
			}
		}
	}

	elapsed := time.Since(start)
	rps := float64(totalOrders) / elapsed.Seconds()

	fmt.Printf("\n✨ LOAD TEST COMPLETED\n")
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("✅ Successful Orders (Inserted in DB): %d\n", successCount)
	fmt.Printf("❌ Failed Requests:                    %d\n", failCount)
	fmt.Printf("⏱️  Total Duration:                     %v\n", elapsed)
	fmt.Printf("🚀 Throughput:                         %.2f Requests/sec\n", rps)
	fmt.Printf("-----------------------------------\n")
}
