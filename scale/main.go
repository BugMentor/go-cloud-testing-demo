package main

import (
	"fmt"
	"sync"
	"time"
)

// Simulates data generation (e.g., external API call or DB query)
func generateData(id int) string {
	time.Sleep(10 * time.Millisecond) // Simulates latency
	return fmt.Sprintf("Record_%d", id)
}

// Worker: Reads jobs from the channel and sends processed data to results
func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		data := generateData(j)
		results <- data
	}
}

func main() {
	const numJobs = 5000
	const numWorkers = 8

	fmt.Printf("Starting processing of %d jobs with %d workers...\n", numJobs, numWorkers)

	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)
	var wg sync.WaitGroup

	start := time.Now()

	// 1. Start Workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Dispatcher: Send jobs to the channel
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// 3. Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect results
	count := 0
	for range results {
		count++
	}

	elapsed := time.Since(start)
	fmt.Printf("Processed %d records in %v\n", count, elapsed)
}
